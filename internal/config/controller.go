package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	_replaceMe                     = "<replace_me>"
	_defaultPhabricatorAPITokenEnv = "PHAB_API_TOKEN"
	_defaultPhabricatorAccessToken = "PHAB_ACCESS_TOKEN"
	_defaultArcanistFilePath       = "~/.arcrc"

	_defaultGithubURL         = "https://github.com/"
	_defaultGithubAPITokenEnv = "GITHUB_API_TOKEN"
	_mode                     = 0o666
)

// Operations is the interface for configuration operations.
type Operations interface {
	Read(path string) (*Config, error)
	Create(path string) error
}

// Controller is the configuration controller
// it handles reading and writing configuration files.
type Controller struct {
	// for testing purposes
	readFile      func(name string) ([]byte, error)
	createFile    func(name string) (*os.File, error)
	mkdirAll      func(path string, perm os.FileMode) error
	yamlUnmarshal func(in []byte, out interface{}) error
	yamMarshal    func(in interface{}) ([]byte, error)
}

var _ Operations = (*Controller)(nil)

// New returns a new configuration controller
func New() *Controller {
	return &Controller{
		readFile:      os.ReadFile,
		createFile:    os.Create,
		mkdirAll:      os.MkdirAll,
		yamlUnmarshal: yaml.Unmarshal,
		yamMarshal:    yaml.Marshal,
	}
}

// Read reads the configuration file at the given path
func (c *Controller) Read(path string) (*Config, error) {
	yamlFile, err := c.readFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file %s: %w", path, err)
	}
	cfg := &Config{}
	if err := c.yamlUnmarshal(yamlFile, cfg); err != nil {
		return nil, fmt.Errorf("unable to parse config.yaml: %w", err)
	}
	return cfg, nil
}

// Create creates a new default configuration file at the given path
func (c *Controller) Create(path string) error {
	// create the file
	if err := c.mkdirAll(filepath.Dir(path), _mode); err != nil {
		return err
	}
	f, err := c.createFile(path)
	if err != nil {
		return err
	}
	defer f.Close()
	// write the default configuration to the file
	out, err := c.yamMarshal(defaultConfiguration())
	if err != nil {
		return err
	}
	_, err = f.Write(out)
	return err
}
