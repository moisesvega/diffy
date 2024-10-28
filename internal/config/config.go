package config

// Config root configuration structure
type Config struct {
	APIs  APIs  `yaml:"apis"`
	Teams Teams `yaml:"teams"`
}

type Me struct {
	PhabricatorUsername string `yaml:"phabricator"`
}

type APIs struct {
	Phabricator Phabricator `yaml:"phabricator"`
}

// Phabricator configuration to set a proper Phabricator Client
type Phabricator struct {
	URL            string `yaml:"base_url"`
	APIToken       string `yaml:"api_token"`
	APITokenEnv    string `yaml:"api_token_env"`
	AccessToken    string `yaml:"access_token"`
	AccessTokenEnv string `yaml:"access_token_env"`
}

type Teams map[string]Team

type Team struct {
	PhabricatorUsers []string `yaml:"phabricator_users"`
	GithubUsers      []string `yaml:"github_users"`
}

// DefaultConfiguration returns default configuration.
func DefaultConfiguration() *Config {
	return &Config{
		APIs: APIs{
			Phabricator: Phabricator{
				URL:            _replaceMe,
				APITokenEnv:    _defaultPhabricatorAPITokenEnv,
				AccessTokenEnv: _defaultPhabricatorAccessToken,
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
