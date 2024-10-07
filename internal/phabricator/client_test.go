package phabricator

import (
	"fmt"
	"testing"

	"github.com/moisesvega/diffy/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/uber/gonduit/test/server"
)

func TestNew(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		s := server.New()
		t.Cleanup(s.Close)
		s.RegisterCapabilities()
		c, err := New(&config.Phabricator{APIToken: "1", URL: s.GetURL(), AccessToken: "1"})
		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("errors", func(t *testing.T) {
		c, err := New(nil)
		require.Error(t, err)
		require.Nil(t, c)
	})
}

func TestClientRequiredConfig(t *testing.T) {
	tests := []struct {
		desc string
		give *config.Phabricator
		want error
	}{
		{desc: "api_token_not_provided", give: &config.Phabricator{}, want: errNoAPITokenProvided},
		{desc: "url_not_provided", give: &config.Phabricator{APIToken: "1"}, want: errNoURLProvided},
		{
			desc: "arcrc_not_found",
			give: &config.Phabricator{APIToken: "1", URL: "1", ArcrcFilePath: "impossible"},
			want: fmt.Errorf("unable to find .arrcrc file"),
		},
		{
			desc: "url_not_provided",
			give: &config.Phabricator{APIToken: "1", URL: "1", ArcrcFilePath: "impossible"},
			want: fmt.Errorf("unable to find .arrcrc file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := New(tt.give)
			require.Error(t, err)
			require.Equal(t, tt.want, err)
		})
	}
}
