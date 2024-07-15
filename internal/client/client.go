package client

import (
	"context"
	"fmt"

	"github.com/uber/gonduit"
	"github.com/uber/gonduit/core"
	"github.com/uber/gonduit/requests"
	"golang.org/x/oauth2"
)

type Client interface {
	GetUsers(strings []string) ([]*User, error)
}

type client struct {
	conn *gonduit.Conn
}

func (c *client) GetUsers(names []string) ([]*User, error) {
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

func New(token string) (Client, error) {
	conn, err := createConnection(token)
	return &client{conn: conn}, err
}

type whoamiresponse struct {
	PHID  string `json:"phid"`
	Email string `json:"email"`
}

const _codeUrl = "https://someurl"

func createConnection(token string) (*gonduit.Conn, error) {
	arcToken, err := getArcToken()
	if err != nil {
		return nil, err
	}
	oc := oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: arcToken}),
	)
	opts := core.ClientOptions{
		Client:   oc,
		APIToken: arcToken,
	}

	conn, err := gonduit.Dial(_codeUrl, &opts)
	if err != nil {
		return nil, err
	}
	whoami := whoamiresponse{}
	if err = conn.Call("user.whoami", &requests.Request{}, &whoami); err != nil {
		return nil, fmt.Errorf("unable to call user.whoami: %w", err)
	}

	return conn, nil
}
