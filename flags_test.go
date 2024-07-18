package main

import (
	"testing"

	"github.com/moisesvega/diffy/internal/config"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFlagSet(t *testing.T) {
	tests := []struct {
		desc string
		give []string
		want *config.Config
	}{
		{desc: "no flags", give: []string{}, want: &config.Config{}},
		{desc: "settings", give: []string{"--settings"}, want: &config.Config{Settings: true}},
		{desc: "phab_users", give: []string{"--phab_users=user,user1"}, want: &config.Config{PhabUsers: []string{"user", "user1"}}},
		{desc: "github_users", give: []string{"--github_users=user,user1"}, want: &config.Config{GithubUsers: []string{"user", "user1"}}},
		// Phab Config
		{desc: "phab_url", give: []string{"--phab_url=someString"}, want: &config.Config{PhabricatorConfig: config.PhabricatorConfig{URL: "someString"}}},
		{desc: "phab_api_token", give: []string{"--phab_api_token=someString"}, want: &config.Config{PhabricatorConfig: config.PhabricatorConfig{APIToken: "someString"}}},
		{desc: "phab_access_token", give: []string{"--phab_access_token=someString"}, want: &config.Config{PhabricatorConfig: config.PhabricatorConfig{AccessToken: "someString"}}},
		{desc: "arrc_file", give: []string{"--arrc_file=someString"}, want: &config.Config{PhabricatorConfig: config.PhabricatorConfig{ArcrcFilePath: "someString"}}},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cfg := &config.Config{}
			pfs := &pflag.FlagSet{}
			fs := registerFlags(pfs, cfg)
			err := fs.Parse(tt.give)
			require.NoError(t, err)
			require.NotNil(t, cfg)
			assert.Equal(t, tt.want, cfg)
		})
	}
}
