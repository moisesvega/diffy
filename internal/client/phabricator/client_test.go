package phabricator

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/moisesvega/diffy/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/uber/gonduit"
	"github.com/uber/gonduit/entities"
	"github.com/uber/gonduit/responses"
	"github.com/uber/gonduit/test/server"
)

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := server.New()
		t.Cleanup(s.Close)
		s.RegisterCapabilities()
		c, err := New(config.Phabricator{APIToken: "1", URL: s.GetURL(), AccessToken: "1"})
		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("errors", func(t *testing.T) {
		c, err := New(config.Phabricator{})
		require.Error(t, err)
		require.NotNil(t, c)
	})
}

func TestClientRequiredConfig(t *testing.T) {
	tests := []struct {
		desc string
		give config.Phabricator
		want error
	}{
		{
			desc: "api_token_not_provided",
			give: config.Phabricator{},
			want: errNoAPITokenProvided,
		},
		{
			desc: "url_not_provided",
			give: config.Phabricator{APIToken: "1"},
			want: errNoURLProvided,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := New(tt.give)
			require.Error(t, err)
			require.ErrorIs(t, err, tt.want)
		})
	}
}

// _userQueryMethod is the method name on Phabricator API.
// source: https://github.com/uber/gonduit/blob/master/user.go
// TODO(moises): This is not exposed in the gonduit package, should we expose it?
const _userQueryMethod = "user.query"

func TestClient_GetUsers(t *testing.T) {
	// CreateDefaults a test server
	s := server.New()
	// Close the server when the test finishes
	t.Cleanup(s.Close)

	user := "someuser"
	userPHID := "PHID-USER-123"

	uqResponse, err := json.Marshal(resultWrapper{
		responses.UserQueryResponse{
			entities.User{
				UserName: user,
				PHID:     userPHID,
			},
		}})
	require.NoError(t, err)

	// First call to the server to get the user
	// Register the method, the status code and the response
	s.RegisterMethod(_userQueryMethod, http.StatusOK, server.ResponseFromJSON(string(uqResponse)))
	authorDifferentials, err := json.Marshal(resultWrapper{
		responses.DifferentialQueryResponse{
			&entities.DifferentialRevision{
				AuthorPHID: userPHID,
			},
		}})
	require.NoError(t, err)

	// In the second call we will call the server to get the differentials
	// where the user is the author
	s.RegisterMethod(gonduit.DifferentialQueryMethod, http.StatusOK, server.ResponseFromJSON(string(authorDifferentials)))

	reviewerDifferential, err := json.Marshal(resultWrapper{
		responses.DifferentialQueryResponse{
			&entities.DifferentialRevision{
				AuthorPHID: userPHID,
			},
		}})
	require.NoError(t, err)
	// In the third call we will call the server to get the differentials
	// where the user is the reviewer
	s.RegisterMethod(gonduit.DifferentialQueryMethod, http.StatusOK, server.ResponseFromJSON(string(reviewerDifferential)))

	// Register the capabilities
	// This is used to register the methods that the server will respond to
	s.RegisterCapabilities()

	// Second call to the server to get the diffs
	// Register the method, the status code and the response
	// s.RegisterMethod(gonduit.DifferentialQueryMethod, http.StatusOK, server.ResponseFromJSON(string(diffs)))
	c, err := New(config.Phabricator{APIToken: "1", URL: s.GetURL(), AccessToken: "1"})
	require.NoError(t, err)
	got, err := c.GetUsers([]string{user})
	require.NoError(t, err)
	require.NotNil(t, got)
}

// resultWrapper is a wrapper for the response from the server
// this is needed because the server returns a response with a "result" key
// source: https://github.com/uber/gonduit/blob/master/core/call.go#L65
type resultWrapper struct {
	Result interface{} `json:"result"`
}
