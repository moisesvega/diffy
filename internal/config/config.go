package config

type Config struct {
	Settings          bool
	PhabUsers         []string
	GithubUsers       []string
	PhabricatorConfig PhabricatorConfig
}

type PhabricatorConfig struct {
	ArcrcFilePath string `yaml:"arrc_file_path"`
	APIToken      string `yaml:"api_token"`
	AccessToken   string `yaml:"access_token"`
}
