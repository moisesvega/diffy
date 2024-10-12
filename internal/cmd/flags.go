package cmd

import (
	"github.com/moisesvega/diffy/internal/config"
	"github.com/spf13/pflag"
)

func setFlags(fs *pflag.FlagSet, cfg *config.Config) {
	// TODO: Create proper descriptions
	fs.BoolVar(&cfg.Settings, "settings", false, "edit settings")
	fs.StringSliceVar(&cfg.PhabUsers, "phab_users", nil, "List of phabricator users you want to track.")
	fs.StringSliceVar(&cfg.GithubUsers, "github_users", nil, "List of github users you want to track.")

	fs.SortFlags = false
}
