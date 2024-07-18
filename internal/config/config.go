package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config root configuration structure
type Config struct {
	Settings    bool
	PhabUsers   []string `yaml:"phab_users"`
	GithubUsers []string `yaml:"github_users"`

	PhabricatorConfig PhabricatorConfig `yaml:"phabricator_config"`
}

// PhabricatorConfig configuration to set a proper Phabricator Client
type PhabricatorConfig struct {
	URL            string `yaml:"url"`
	APIToken       string `yaml:"api_token"`
	APITokenEnv    string `yaml:"api_token_env"`
	AccessToken    string `yaml:"access_token"`
	AccessTokenEnv string `yaml:"access_token_env"`
	ArcrcFilePath  string `yaml:"arrc_file_path"`
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

// DefaultConfiguration returns default configuration.
func DefaultConfiguration() *Config {
	return &Config{
		PhabricatorConfig: PhabricatorConfig{
			APITokenEnv:    "PHAB_API_TOKEN",
			AccessTokenEnv: "PHAB_ACCESS_TOKEN",
			ArcrcFilePath:  "~/.arcrc",
		},
	}
}
