package phabricator

import (
	"context"
	"errors"
	"os"

	"github.com/moisesvega/diffy/internal/config"
	"github.com/uber/gonduit"
	"github.com/uber/gonduit/core"
	"github.com/uber/gonduit/requests"
	"github.com/uber/gonduit/responses"
	"golang.org/x/oauth2"
)

//go:generate mockgen -source=client.go -destination=mock_phabricator/mocks.go -self_package=github.com/moisesvega/diffy/internal/client/phabricator/mock_phabricator

type PhabClient interface {
	GetUsers(strings []string) ([]*User, error)
}

type Client struct {
	conn igonduit
}

type igonduit interface {
	UserQuery(req requests.UserQueryRequest) (*responses.UserQueryResponse, error)
	DifferentialQuery(req requests.DifferentialQueryRequest) (*responses.DifferentialQueryResponse, error)
}

var (
	errNoAPITokenProvided = errors.New("no API token provided")
	errNoURLProvided      = errors.New("no URL provided")
)

func New(cfg *config.PhabricatorConfig) (PhabClient, error) {
	if cfg == nil {
		return nil, errors.New("phabricator config is required")
	}
	conn, err := createConnection(cfg)
	return &Client{conn: conn}, err
}

func createConnection(cfg *config.PhabricatorConfig) (*gonduit.Conn, error) {
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

func (c *Client) GetUsers(names []string) ([]*User, error) {
	res, err := c.conn.UserQuery(requests.UserQueryRequest{
		Usernames: names,
	})
	if err != nil {
		return nil, err
	}
	users := make([]*User, len(*res))
	for _, user := range *res {
		users = append(users, &User{User: user})
	}
	for _, u := range users {
		diffs, err := c.conn.DifferentialQuery(requests.DifferentialQueryRequest{
			Authors: []string{u.PHID},
		})
		if err != nil {
			return nil, err
		}
		u.Differentials = *diffs

		reviews, err := c.conn.DifferentialQuery(requests.DifferentialQueryRequest{
			Reviewers: []string{u.PHID},
		})
		if err != nil {
			return nil, err
		}
		u.Reviews = *reviews
	}
	return users, nil
}

// type whoamiresponse struct {
// 	PHID  string `json:"phid"`
// 	Email string `json:"email"`
// }
//
// func (c *Client) CheckConnection() error {
// 	whoami := whoamiresponse{}
// 	if err := c.conn.Call("user.whoami", &requests.Request{}, &whoami); err != nil {
// 		return fmt.Errorf("unable to call user.whoami: %w", err)
// 	}
// 	return nil
// }
