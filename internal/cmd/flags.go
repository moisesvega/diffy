package cmd

import (
	"github.com/spf13/pflag"
)

type opts struct {
	settings  bool
	phabUsers []string
}

const _settingsDesc = `Opens the settings file in the editor. By default, uses $XDG_CONFIG_HOME as the path. On macOS, if $XDG_CONFIG_HOME is not set, defaults to $HOME/Library/Application Support/diffy`

func setFlags(fs *pflag.FlagSet, opts *opts) {
	fs.BoolVar(&opts.settings, "settings", false, _settingsDesc)
	fs.StringSliceVar(&opts.phabUsers, "phab_users", nil, "List of phabricator users you want to track.")
	// TODO(moisesvega): Add github users flag

	fs.SortFlags = false
}
