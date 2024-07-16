package main

import (
	"github.com/moisesvega/diffy/internal/config"
	"github.com/spf13/pflag"
)

func registerFlags(fs *pflag.FlagSet, cfg *config.Config) *pflag.FlagSet {
	fs.BoolVar(&cfg.Settings, "settings", false, "edit settings")
	fs.StringSliceVar(&cfg.PhabUsers, "phab_users", nil, "phab_users=user_1,user_2")
	fs.StringSliceVar(&cfg.GithubUsers, "github_users", nil, "github_users=user_1,user_2")

	// Phabricator Configuration
	fs.StringVar(&cfg.PhabricatorConfig.APIToken, "api_token", "", "api_token=token")
	fs.StringVar(&cfg.PhabricatorConfig.AccessToken, "access_token", "", "access_token=token")
	fs.StringVar(&cfg.PhabricatorConfig.ArcrcFilePath, "arrc_file", "~/.arcrc", "arrc_file=~/.arcrc")

	// Github Configuration
	return fs
}
