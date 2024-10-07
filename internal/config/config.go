package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config root configuration structure
type Config struct {
	Me    Me    `yaml:"me"`
	APIs  APIs  `yaml:"apis"`
	Teams Teams `yaml:"teams"`
	// Use for the CLI
	Settings    bool
	PhabUsers   []string
	GithubUsers []string
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

// ReadConfiguration from a given filepath
func ReadConfiguration(filepath string) (*Config, error) {
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file %s: %w", filepath, err)
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(yamlFile, cfg); err != nil {
		return nil, fmt.Errorf("unable to parse config.yaml: %w", err)
	}
	return cfg, nil
}

const (
	_replaceMe                     = "<replace_me>"
	_defaultPhabricatorAPITokenEnv = "PHAB_API_TOKEN"
	_defaultPhabricatorAccessToken = "PHAB_ACCESS_TOKEN"
	_defaultArcanistFilePath       = "~/.arcrc"

	_defaultGithubURL         = "https://github.com/"
	_defaultGithubAPITokenEnv = "GITHUB_API_TOKEN"
)

// DefaultConfiguration returns default configuration.
func DefaultConfiguration() *Config {
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
