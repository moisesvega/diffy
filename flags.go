package main

import (
	"github.com/moisesvega/diffy/internal/config"
	"github.com/spf13/pflag"
)

func registerFlags(fs *pflag.FlagSet, cfg *config.Config) *pflag.FlagSet {
	fs.BoolVar(&cfg.Settings, "settings", false, "edit settings")
	fs.StringArrayVar(&cfg.PhabUsers, "phab_users", make([]string, 0), "phab_users=user_1,user_2")
	fs.StringArrayVar(&cfg.GithubUsers, "github_users", make([]string, 0), "github_users=user_1,user_2")
	return fs
}
