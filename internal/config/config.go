package config

// Config root configuration structure
type Config struct {
	Me    Me    `yaml:"me"`
	APIs  APIs  `yaml:"apis"`
	Teams Teams `yaml:"teams"`
}

type Me struct {
	PhabricatorUsername string `yaml:"phabricator"`
	GithubUsername      string `yaml:"github"`
}

type APIs struct {
	Phabricator Phabricator `yaml:"phabricator"`
	Github      Github      `yaml:"github"`
}

// Phabricator configuration to set a proper Phabricator Client
type Phabricator struct {
	URL            string `yaml:"base_url"`
	APIToken       string `yaml:"api_token"`
	APITokenEnv    string `yaml:"api_token_env"`
	AccessToken    string `yaml:"access_token"`
	AccessTokenEnv string `yaml:"access_token_env"`
	ArcrcFilePath  string `yaml:"arrc_file_path"`
}

type Github struct {
	URL         string `yaml:"base_url"`
	APIToken    string `yaml:"api_token"`
	APITokenEnv string `yaml:"api_token_env"`
}

type Teams map[string]Team

type Team struct {
	PhabricatorUsers []string `yaml:"phabricator_users"`
	GithubUsers      []string `yaml:"github_users"`
}

// defaultConfiguration returns default configuration.
func defaultConfiguration() *Config {
	return &Config{
		Me: Me{
			PhabricatorUsername: _replaceMe,
			GithubUsername:      _replaceMe,
		},
		APIs: APIs{
			Phabricator: Phabricator{
				URL:            _replaceMe,
				APITokenEnv:    _defaultPhabricatorAPITokenEnv,
				AccessTokenEnv: _defaultPhabricatorAccessToken,
				ArcrcFilePath:  _defaultArcanistFilePath,
			},
			Github: Github{
				URL:         _defaultGithubURL,
				APIToken:    "",
				APITokenEnv: _defaultGithubAPITokenEnv,
			},
		},
		Teams: Teams{
			"a_team": Team{
				PhabricatorUsers: []string{
					_replaceMe,
				},
				GithubUsers: []string{
					_replaceMe,
				},
			},
		},
	}
}
