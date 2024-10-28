package config

// Config root configuration structure
type Config struct {
	APIs APIs `yaml:"apis"`
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
	}
}
