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
	url := "https://example.com"
	somePath := "/home/some/path"
	give := []string{
		"--settings",
		// Users
		"--phab_users=" + usersString,
		"--github_users=" + usersString,

		//  Phab configuration
		"--phab_url=" + url,
		"--phab_api_token=" + token,
		"--phab_access_token=" + token,
		"--arrc_file=" + somePath,
	}

	want := &config.Config{
		Settings:    true,
		PhabUsers:   users,
		GithubUsers: users,
		PhabricatorConfig: config.PhabricatorConfig{
			URL:           url,
			APIToken:      token,
			AccessToken:   token,
			ArcrcFilePath: somePath,
		},
	}

	got := &config.Config{}
	pfs := pflag.NewFlagSet("new", pflag.ExitOnError)
	require.NotPanics(t, func() {
		setFlags(pfs, got)
	})
	err := pfs.Parse(give)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, want, got)
}
