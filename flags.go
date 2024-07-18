package main

import (
	"github.com/moisesvega/diffy/internal/config"
	"github.com/spf13/pflag"
)

func registerFlags(fs *pflag.FlagSet, cfg *config.Config) *pflag.FlagSet {
	// TODO: Create proper descriptions
	fs.BoolVar(&cfg.Settings, "settings", false, "edit settings")
	fs.StringSliceVar(&cfg.PhabUsers, "phab_users", nil, "phab_users=user_1,user_2")
	fs.StringSliceVar(&cfg.GithubUsers, "github_users", nil, "github_users=user_1,user_2")

	// Phabricator Configuration
	// TODO: Look into how to parse the URL in this step
	fs.StringVar(&cfg.PhabricatorConfig.URL, "phab_url", "", "api_token=token")
	fs.StringVar(&cfg.PhabricatorConfig.APIToken, "phab_api_token", "", "api_token=token")
	fs.StringVar(&cfg.PhabricatorConfig.AccessToken, "phab_access_token", "", "access_token=token")
	fs.StringVar(&cfg.PhabricatorConfig.ArcrcFilePath, "arrc_file", "", "arrc_file=~/.arcrc")

	fs.SortFlags = false
	// Github Configuration
	return fs
}
