package github

import (
	"github.com/shurcooL/githubv4"
)

type Client struct {
}

func New() (*Client, error) {

	_ = githubv4.NewClient(nil)

	return nil, nil
}
