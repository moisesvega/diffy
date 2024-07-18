package phabricator

import (
	"fmt"
	"testing"

	"github.com/moisesvega/diffy/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/uber/gonduit/test/server"
)

func TestNewClient(t *testing.T) {
	c := NewClient()
	require.NotNil(t, c)
}

func TestNew(t *testing.T) {
	c := NewClient()

	s := server.New()
	t.Cleanup(s.Close)

	s.RegisterCapabilities()
	require.NotNil(t, c)
	c, err := c.New(&config.PhabricatorConfig{APIToken: "1", URL: s.GetURL(), AccessToken: "1"})
	require.NoError(t, err)
	require.NotNil(t, c)

	t.Run("errors", func(t *testing.T) {
		c, err := c.New(nil)
		require.Error(t, err)
		require.Nil(t, c)
	})
}

func TestClientRequiredConfig(t *testing.T) {
	c := NewClient()
	tests := []struct {
		desc string
		give *config.PhabricatorConfig
		want error
	}{
		{desc: "api_token_not_provided", give: &config.PhabricatorConfig{}, want: errNoAPITokenProvided},
		{desc: "url_not_provided", give: &config.PhabricatorConfig{APIToken: "1"}, want: errNoURLProvided},
		{
			desc: "arcrc_not_found",
			give: &config.PhabricatorConfig{APIToken: "1", URL: "1", ArcrcFilePath: "imposible"},
			want: fmt.Errorf("unable to find .arrcrc file"),
		},
		{
			desc: "url_not_provided",
			give: &config.PhabricatorConfig{APIToken: "1", URL: "1", ArcrcFilePath: "imposible"},
			want: fmt.Errorf("unable to find .arrcrc file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := c.New(tt.give)
			require.Error(t, err)
			require.Equal(t, tt.want, err)
		})
	}
}

//go:generate mockgen -source=client.go -destination=mock_phabricator/mocks.go -self_package=github.com/moisesvega/diffy/internal/client/phabricator/mock_phabricator

func TestClient_GetUsers(t *testing.T) {
	// TODO: Create tests for Get Users:
}
