package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/entity"
)

const (
	_defaultBaseURL   = "https://api.github.com"
	_perPage          = 100
	_apiVersionHeader = "2022-11-28"
)

// Client is the interface for the GitHub client.
type Client interface {
	GetUsers(names []string) ([]*entity.User, error)
}

type client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

var errNoTokenProvided = errors.New("no GitHub token provided")

// New creates a new GitHub client.
func New(cfg config.GitHub) (Client, error) {
	token := os.Getenv(cfg.TokenEnv)
	if cfg.Token != "" {
		token = cfg.Token
	}
	if token == "" {
		return nil, errNoTokenProvided
	}

	baseURL := _defaultBaseURL
	if cfg.BaseURL != "" {
		baseURL = cfg.BaseURL
	}

	return &client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{},
	}, nil
}

// GetUsers returns a list of users with their pull requests and reviews.
func (c *client) GetUsers(names []string) ([]*entity.User, error) {
	users := make([]*entity.User, 0, len(names))
	for _, name := range names {
		user := &entity.User{Username: name}

		authored, err := c.getPRs("author", name)
		if err != nil {
			return nil, err
		}
		user.Differentials = authored

		reviewed, err := c.getPRs("reviewed-by", name)
		if err != nil {
			return nil, err
		}
		user.Reviews = reviewed

		users = append(users, user)
	}
	return users, nil
}

func (c *client) getPRs(qualifier, username string) ([]*entity.Differential, error) {
	query := fmt.Sprintf("type:pr %s:%s", qualifier, username)
	items, err := c.searchIssues(query)
	if err != nil {
		return nil, err
	}

	diffs := make([]*entity.Differential, 0, len(items))
	for _, item := range items {
		if item.PullRequest == nil || item.PullRequest.URL == "" {
			continue
		}
		detail, err := c.getPRDetail(item.PullRequest.URL)
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, mapPRToDifferential(detail))
	}
	return diffs, nil
}

func (c *client) searchIssues(query string) ([]searchItem, error) {
	var allItems []searchItem
	page := 1

	for {
		u := fmt.Sprintf("%s/search/issues?q=%s&per_page=%d&page=%d",
			c.baseURL, url.QueryEscape(query), _perPage, page)

		items, total, err := c.doSearch(u)
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, items...)

		if len(allItems) >= total || len(items) < _perPage {
			break
		}
		page++
	}

	return allItems, nil
}

func (c *client) doSearch(u string) ([]searchItem, int, error) {
	resp, err := c.doRequest(u)
	if err != nil {
		return nil, 0, fmt.Errorf("search request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("search API returned status %d", resp.StatusCode)
	}

	var result searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode search response: %w", err)
	}
	return result.Items, result.TotalCount, nil
}

func (c *client) getPRDetail(prURL string) (*pullRequestDetail, error) {
	resp, err := c.doRequest(prURL)
	if err != nil {
		return nil, fmt.Errorf("PR detail request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PR detail API returned status %d", resp.StatusCode)
	}

	var detail pullRequestDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, fmt.Errorf("failed to decode PR detail: %w", err)
	}
	return &detail, nil
}

func (c *client) doRequest(u string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", _apiVersionHeader)
	return c.httpClient.Do(req)
}

func mapPRToDifferential(pr *pullRequestDetail) *entity.Differential {
	status := entity.Unknown
	if pr.Merged {
		status = entity.Closed
	} else if pr.State == "open" {
		status = entity.Accepted
	}

	return &entity.Differential{
		ID:         strconv.Itoa(pr.Number),
		Title:      pr.Title,
		LineCount:  pr.Additions + pr.Deletions,
		Status:     status,
		StatusName: pr.State,
		URI:        pr.HTMLURL,
		CreatedAt:  pr.CreatedAt,
		ModifiedAt: pr.UpdatedAt,
	}
}

// searchResponse represents the GitHub search API response.
type searchResponse struct {
	TotalCount int          `json:"total_count"`
	Items      []searchItem `json:"items"`
}

type searchItem struct {
	Number        int             `json:"number"`
	Title         string          `json:"title"`
	State         string          `json:"state"`
	HTMLURL       string          `json:"html_url"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	RepositoryURL string          `json:"repository_url"`
	PullRequest   *pullRequestRef `json:"pull_request"`
}

type pullRequestRef struct {
	URL      string     `json:"url"`
	MergedAt *time.Time `json:"merged_at"`
}

// pullRequestDetail represents a GitHub pull request with full details.
type pullRequestDetail struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	Merged    bool      `json:"merged"`
	Additions int       `json:"additions"`
	Deletions int       `json:"deletions"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
