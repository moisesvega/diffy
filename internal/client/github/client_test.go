package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("success with token", func(t *testing.T) {
		c, err := New(config.GitHub{Token: "tok"})
		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("success with token env", func(t *testing.T) {
		t.Setenv("TEST_GH_TOKEN", "tok")
		c, err := New(config.GitHub{TokenEnv: "TEST_GH_TOKEN"})
		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("no token", func(t *testing.T) {
		c, err := New(config.GitHub{})
		require.ErrorIs(t, err, errNoTokenProvided)
		require.Nil(t, c)
	})

	t.Run("custom base url", func(t *testing.T) {
		c, err := New(config.GitHub{Token: "tok", BaseURL: "https://ghe.example.com/api/v3"})
		require.NoError(t, err)
		require.NotNil(t, c)
	})
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	var srv *httptest.Server
	mux := http.NewServeMux()

	mux.HandleFunc("/search/issues", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		var resp searchResponse

		switch {
		case strings.Contains(q, "author:testuser"):
			resp = searchResponse{
				TotalCount: 1,
				Items: []searchItem{
					{
						Number:  1,
						Title:   "Fix bug",
						State:   "closed",
						HTMLURL: "https://github.com/org/repo/pull/1",
						PullRequest: &pullRequestRef{
							URL:      srv.URL + "/repos/org/repo/pulls/1",
							MergedAt: &now,
						},
					},
				},
			}
		case strings.Contains(q, "reviewed-by:testuser"):
			resp = searchResponse{
				TotalCount: 1,
				Items: []searchItem{
					{
						Number:  2,
						Title:   "Add feature",
						State:   "open",
						HTMLURL: "https://github.com/org/repo/pull/2",
						PullRequest: &pullRequestRef{
							URL: srv.URL + "/repos/org/repo/pulls/2",
						},
					},
				},
			}
		default:
			resp = searchResponse{TotalCount: 0}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/repos/org/repo/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(pullRequestDetail{
			Number:    1,
			Title:     "Fix bug",
			State:     "closed",
			Merged:    true,
			Additions: 10,
			Deletions: 5,
			HTMLURL:   "https://github.com/org/repo/pull/1",
			CreatedAt: now,
			UpdatedAt: now,
		})
	})

	mux.HandleFunc("/repos/org/repo/pulls/2", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(pullRequestDetail{
			Number:    2,
			Title:     "Add feature",
			State:     "open",
			Merged:    false,
			Additions: 50,
			Deletions: 3,
			HTMLURL:   "https://github.com/org/repo/pull/2",
			CreatedAt: now,
			UpdatedAt: now,
		})
	})

	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestGetUsers(t *testing.T) {
	srv := newTestServer(t)

	c := &client{
		baseURL:    srv.URL,
		token:      "test-token",
		httpClient: srv.Client(),
	}

	t.Run("success", func(t *testing.T) {
		users, err := c.GetUsers([]string{"testuser"})
		require.NoError(t, err)
		require.Len(t, users, 1)

		user := users[0]
		assert.Equal(t, "testuser", user.Username)

		// Authored PRs
		require.Len(t, user.Differentials, 1)
		assert.Equal(t, "1", user.Differentials[0].ID)
		assert.Equal(t, "Fix bug", user.Differentials[0].Title)
		assert.Equal(t, 15, user.Differentials[0].LineCount)
		assert.Equal(t, entity.Closed, user.Differentials[0].Status)

		// Reviewed PRs
		require.Len(t, user.Reviews, 1)
		assert.Equal(t, "2", user.Reviews[0].ID)
		assert.Equal(t, "Add feature", user.Reviews[0].Title)
		assert.Equal(t, 53, user.Reviews[0].LineCount)
		assert.Equal(t, entity.Accepted, user.Reviews[0].Status)
	})
}

func TestGetUsers_SearchError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	c := &client{
		baseURL:    srv.URL,
		token:      "test-token",
		httpClient: srv.Client(),
	}

	users, err := c.GetUsers([]string{"testuser"})
	require.Error(t, err)
	require.Nil(t, users)
}

func TestGetUsers_PRDetailError(t *testing.T) {
	var srv *httptest.Server
	mux := http.NewServeMux()

	mux.HandleFunc("/search/issues", func(w http.ResponseWriter, r *http.Request) {
		resp := searchResponse{
			TotalCount: 1,
			Items: []searchItem{
				{
					Number: 1,
					PullRequest: &pullRequestRef{
						URL: srv.URL + "/repos/org/repo/pulls/1",
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/repos/org/repo/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	c := &client{
		baseURL:    srv.URL,
		token:      "test-token",
		httpClient: srv.Client(),
	}

	users, err := c.GetUsers([]string{"testuser"})
	require.Error(t, err)
	require.Nil(t, users)
}

func TestMapPRToDifferential(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		pr     *pullRequestDetail
		status entity.Status
	}{
		{
			name: "merged PR maps to Closed",
			pr: &pullRequestDetail{
				Number: 1, Title: "PR", State: "closed", Merged: true,
				Additions: 10, Deletions: 5, HTMLURL: "url",
				CreatedAt: now, UpdatedAt: now,
			},
			status: entity.Closed,
		},
		{
			name: "open PR maps to Accepted",
			pr: &pullRequestDetail{
				Number: 2, Title: "PR", State: "open", Merged: false,
				Additions: 20, Deletions: 0, HTMLURL: "url",
				CreatedAt: now, UpdatedAt: now,
			},
			status: entity.Accepted,
		},
		{
			name: "closed unmerged PR maps to Unknown",
			pr: &pullRequestDetail{
				Number: 3, Title: "PR", State: "closed", Merged: false,
				Additions: 5, Deletions: 5, HTMLURL: "url",
				CreatedAt: now, UpdatedAt: now,
			},
			status: entity.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mapPRToDifferential(tt.pr)
			assert.Equal(t, tt.status, d.Status)
			assert.Equal(t, tt.pr.Additions+tt.pr.Deletions, d.LineCount)
		})
	}
}

func TestDoRequest_AuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/vnd.github+json", r.Header.Get("Accept"))
		assert.Equal(t, _apiVersionHeader, r.Header.Get("X-GitHub-Api-Version"))
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	c := &client{
		baseURL:    srv.URL,
		token:      "test-token",
		httpClient: srv.Client(),
	}

	resp, err := c.doRequest(srv.URL + "/test")
	require.NoError(t, err)
	_ = resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSearchPagination(t *testing.T) {
	callCount := 0
	var srv *httptest.Server
	mux := http.NewServeMux()

	mux.HandleFunc("/search/issues", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		page := r.URL.Query().Get("page")

		var resp searchResponse
		if page == "" || page == "1" {
			items := make([]searchItem, _perPage)
			for i := range items {
				items[i] = searchItem{
					Number: i + 1,
					PullRequest: &pullRequestRef{
						URL: srv.URL + "/repos/org/repo/pulls/1",
					},
				}
			}
			resp = searchResponse{TotalCount: _perPage + 1, Items: items}
		} else {
			resp = searchResponse{
				TotalCount: _perPage + 1,
				Items: []searchItem{
					{
						Number: _perPage + 1,
						PullRequest: &pullRequestRef{
							URL: srv.URL + "/repos/org/repo/pulls/1",
						},
					},
				},
			}
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/repos/org/repo/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(pullRequestDetail{
			Number: 1, Title: "PR", State: "closed", Merged: true,
			Additions: 1, Deletions: 0, HTMLURL: "url",
		})
	})

	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	c := &client{
		baseURL:    srv.URL,
		token:      "test-token",
		httpClient: srv.Client(),
	}

	diffs, err := c.getPRs("author", "user")
	require.NoError(t, err)
	assert.Len(t, diffs, _perPage+1)
	assert.Equal(t, 2, callCount)
}

func TestDoSearch_DecodeError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("invalid json"))
	}))
	t.Cleanup(srv.Close)

	c := &client{baseURL: srv.URL, token: "tok", httpClient: srv.Client()}
	_, _, err := c.doSearch(srv.URL + "/search/issues?q=test")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode search response")
}

func TestGetPRDetail_DecodeError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not json"))
	}))
	t.Cleanup(srv.Close)

	c := &client{baseURL: srv.URL, token: "tok", httpClient: srv.Client()}
	_, err := c.getPRDetail(srv.URL + "/repos/org/repo/pulls/1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode PR detail")
}

func TestDoRequest_InvalidURL(t *testing.T) {
	c := &client{token: "tok", httpClient: &http.Client{}}
	_, err := c.doRequest("://invalid")
	require.Error(t, err)
}

func TestGetPRs_NilPullRequest(t *testing.T) {
	var srv *httptest.Server
	mux := http.NewServeMux()

	mux.HandleFunc("/search/issues", func(w http.ResponseWriter, r *http.Request) {
		resp := searchResponse{
			TotalCount: 2,
			Items: []searchItem{
				{Number: 1, PullRequest: nil},
				{Number: 2, PullRequest: &pullRequestRef{URL: ""}},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	c := &client{baseURL: srv.URL, token: "tok", httpClient: srv.Client()}
	diffs, err := c.getPRs("author", "user")
	require.NoError(t, err)
	assert.Empty(t, diffs)
}

func TestGetUsers_ReviewError(t *testing.T) {
	var srv *httptest.Server
	mux := http.NewServeMux()
	callCount := 0

	mux.HandleFunc("/search/issues", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			// First call (authored PRs) succeeds with empty results
			_ = json.NewEncoder(w).Encode(searchResponse{TotalCount: 0})
		} else {
			// Second call (reviewed PRs) fails
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	c := &client{baseURL: srv.URL, token: "tok", httpClient: srv.Client()}
	users, err := c.GetUsers([]string{"testuser"})
	require.Error(t, err)
	require.Nil(t, users)
}

func TestDoSearch_RequestError(t *testing.T) {
	c := &client{baseURL: "http://localhost:0", token: "tok", httpClient: &http.Client{}}
	_, _, err := c.doSearch("http://localhost:0/search/issues?q=test")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "search request failed")
}

func TestGetPRDetail_RequestError(t *testing.T) {
	c := &client{baseURL: "http://localhost:0", token: "tok", httpClient: &http.Client{}}
	_, err := c.getPRDetail("http://localhost:0/repos/org/repo/pulls/1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "PR detail request failed")
}
