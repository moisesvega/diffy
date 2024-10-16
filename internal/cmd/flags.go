package cmd

import (
	"github.com/spf13/pflag"
)

type opts struct {
	Settings    bool
	PhabUsers   []string
	GithubUsers []string
}

func setFlags(fs *pflag.FlagSet, opts *opts) {
	// TODO: Create proper descriptions
	fs.BoolVar(&opts.Settings, "settings", false, "edit settings")
	fs.StringSliceVar(&opts.PhabUsers, "phab_users", nil, "List of phabricator users you want to track.")
	fs.StringSliceVar(&opts.GithubUsers, "github_users", nil, "List of github users you want to track.")

	fs.SortFlags = false
}
