package cmd

import (
	"github.com/spf13/pflag"
)

type opts struct {
	settings    bool
	phabUsers   []string
	githubUsers []string
}

func setFlags(fs *pflag.FlagSet, opts *opts) {
	// TODO: CreateDefaults proper descriptions
	fs.BoolVar(&opts.settings, "settings", false, "edit settings")
	fs.StringSliceVar(&opts.phabUsers, "phab_users", nil, "List of phabricator users you want to track.")
	fs.StringSliceVar(&opts.githubUsers, "github_users", nil, "List of github users you want to track.")

	fs.SortFlags = false
}
