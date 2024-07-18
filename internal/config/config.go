package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Settings          bool
	PhabUsers         []string          `yaml:"phab_users"`
	GithubUsers       []string          `yaml:"github_users"`
	PhabricatorConfig PhabricatorConfig `yaml:"phabricator_config"`
}

type PhabricatorConfig struct {
	URL           string `yaml:"url"`
	APIToken      string `yaml:"api_token"`
	AccessToken   string `yaml:"access_token"`
	ArcrcFilePath string `yaml:"arrc_file_path"`
}

func ReadConfiguration(filepath string) (*Config, error) {
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file %s: %w", filepath, err)
	}
	fmt.Println(string(yamlFile))
	cfg := &Config{}
	if err := yaml.Unmarshal(yamlFile, cfg); err != nil {
		return nil, fmt.Errorf("unable to parse config.yaml: %w", err)
	}
	return cfg, nil
}
