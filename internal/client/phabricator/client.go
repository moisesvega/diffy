package phabricator

import (
	"context"
	"errors"
	"os"

	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/mapper"
	"github.com/moisesvega/diffy/internal/model"
	"github.com/uber/gonduit"
	"github.com/uber/gonduit/constants"
	"github.com/uber/gonduit/core"
	"github.com/uber/gonduit/requests"
	"golang.org/x/oauth2"
)

type Client struct {
	conn *gonduit.Conn
}

var (
	errNoAPITokenProvided = errors.New("no API token provided")
	errNoURLProvided      = errors.New("no URL provided")
)

// New creates a new Phabricator Client
func New(cfg config.Phabricator) (*Client, error) {
	conn, err := createConnection(cfg)
	return &Client{conn: conn}, err
}

func createConnection(cfg config.Phabricator) (*gonduit.Conn, error) {
	if len(cfg.APIToken) == 0 && len(cfg.APITokenEnv) == 0 && len(cfg.ArcrcFilePath) == 0 {
		return nil, errNoAPITokenProvided
	}

	if len(cfg.URL) == 0 {
		return nil, errNoURLProvided
	}

	accessToken := os.Getenv(cfg.AccessTokenEnv)
	if len(cfg.AccessToken) > 0 {
		accessToken = cfg.AccessToken
	}

	apiToken := os.Getenv(cfg.APITokenEnv)
	if len(cfg.ArcrcFilePath) > 0 {
		var err error
		apiToken, err = getArcToken(cfg.ArcrcFilePath)
		if err != nil {
			return nil, err
		}
	}
	if len(cfg.APIToken) > 0 {
		apiToken = cfg.APIToken
	}

	oc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))
	return gonduit.Dial(cfg.URL, &core.ClientOptions{
		Client:   oc,
		APIToken: apiToken,
	})
}

// GetUsers returns a list of users with their differentials and reviews.
func (c *Client) GetUsers(names []string) ([]*model.User, error) {
	// We can't query for differentials and reviews by username.
	// We need to query for users first and then query for differentials and reviews by user PHID.
	// This is a limitation of the Phabricator API.
	// First we query for users.
	// Then we query for differentials and reviews by user PHID.
	res, err := c.conn.UserQuery(requests.UserQueryRequest{
		Usernames: names,
	})
	if err != nil {
		return nil, err
	}
	users := make([]*model.User, 0, len(*res))
	for _, user := range *res {
		u := mapper.FromPhabricatorUser(user)
		// Then we query for differentials and reviews by user PHID.
		// TODO(moisesvega): If I see degradation in performance, I will consider get all differentials and reviews in one call.
		if u.Differentials, err = c.getDifferentials(user.PHID); err != nil {
			return nil, err
		}
		if u.Reviews, err = c.getReviews(user.PHID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (c *Client) getDifferentials(id string) ([]*model.Differential, error) {
	res, err := c.conn.DifferentialQuery(requests.DifferentialQueryRequest{
		Status: constants.DifferentialStatusLegacyPublished,
		PHIDs:  []string{id},
	})
	if err != nil {
		return nil, err
	}
	return mapper.FromPhabricatorDifferentialQueryResponse(*res), nil
}

func (c *Client) getReviews(id string) ([]*model.Differential, error) {
	res, err := c.conn.DifferentialQuery(requests.DifferentialQueryRequest{
		Reviewers: []string{id},
	})
	if err != nil {
		return nil, err
	}
	return mapper.FromPhabricatorDifferentialQueryResponse(*res), nil
}
