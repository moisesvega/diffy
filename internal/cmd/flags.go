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

	// Phabricator Configuration
	// TODO: Look into how to parse the URL in this step
	fs.StringVar(&cfg.APIs.Phabricator.URL, "phab_url", "", "Phabricator API url")
	fs.StringVar(&cfg.APIs.Phabricator.APIToken, "phab_api_token", "", "API token this could be inside of your .arrc")
	fs.StringVar(&cfg.APIs.Phabricator.AccessToken, "phab_access_token", "", "If API URL is protected by oauth you can provide your access token here")
	fs.StringVar(&cfg.APIs.Phabricator.ArcrcFilePath, "arrc_file", "", "If provided it will read the .arcrc and get URL and API Token")

	// Github Configurati
	// TODO: Create github configuration
	fs.SortFlags = false
}
