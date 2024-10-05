package cmd

import (
	"strings"
	"testing"

	"github.com/moisesvega/diffy/internal/config"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFlagSet(t *testing.T) {
	users := []string{"first", "second"}
	usersString := strings.Join(users, ",")
	token := "someToken"
	tests := []struct {
		desc string
		give []string
		want *config.Config
	}{
		{desc: "no flags", give: []string{}, want: &config.Config{}},
		{desc: "settings", give: []string{"--settings"}, want: &config.Config{Settings: true}},
		{desc: "phab_users", give: []string{"--phab_users=" + usersString}, want: &config.Config{PhabUsers: users}},
		{desc: "github_users", give: []string{"--github_users=" + usersString}, want: &config.Config{GithubUsers: users}},
		// Phab Config
		{desc: "phab_url", give: []string{"--phab_url=" + token}, want: &config.Config{PhabricatorConfig: config.PhabricatorConfig{URL: token}}},
		{desc: "phab_api_token", give: []string{"--phab_api_token=" + token}, want: &config.Config{PhabricatorConfig: config.PhabricatorConfig{APIToken: token}}},
		{desc: "phab_access_token", give: []string{"--phab_access_token=" + token}, want: &config.Config{PhabricatorConfig: config.PhabricatorConfig{AccessToken: token}}},
		{desc: "arrc_file", give: []string{"--arrc_file=" + token}, want: &config.Config{PhabricatorConfig: config.PhabricatorConfig{ArcrcFilePath: token}}},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cfg := &config.Config{}
			pfs := pflag.NewFlagSet("new", pflag.ExitOnError)
			require.NotPanics(t, func() {
				setFlags(pfs, cfg)
			})
			err := pfs.Parse(tt.give)
			require.NoError(t, err)
			require.NotNil(t, cfg)
			assert.Equal(t, tt.want, cfg)
		})
	}
}
