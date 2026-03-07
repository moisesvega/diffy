package config

// Config root configuration structure
type Config struct {
	APIs APIs `yaml:"apis"`
}

type APIs struct {
	Phabricator Phabricator `yaml:"phabricator"`
	GitHub      GitHub      `yaml:"github"`
}

// GitHub configuration to set a proper GitHub Client
type GitHub struct {
	BaseURL  string `yaml:"base_url"`
	Token    string `yaml:"token"`
	TokenEnv string `yaml:"token_env"`
}

// Phabricator configuration to set a proper Phabricator Client
type Phabricator struct {
	URL            string `yaml:"base_url"`
	APIToken       string `yaml:"api_token"`
	APITokenEnv    string `yaml:"api_token_env"`
	AccessToken    string `yaml:"access_token"`
	AccessTokenEnv string `yaml:"access_token_env"`
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
			GitHub: GitHub{
				TokenEnv: _defaultGitHubTokenEnv,
			},
		},
	}
}
