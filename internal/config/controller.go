package config

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

const (
	_replaceMe                     = "<replace_me>"
	_defaultPhabricatorAPITokenEnv = "PHAB_API_TOKEN"
	_defaultPhabricatorAccessToken = "PHAB_ACCESS_TOKEN"
	_mode                          = os.FileMode(0o677)
)

// Operations is the interface for configuration operations.
type Operations interface {
	Read(path string) (*Config, error)
	CreateDefaults(path string) error
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

// CreateDefaults creates a new default configuration file at the given path
func (c *Controller) CreateDefaults(path string) (err error) {
	// create the file
	if err = c.mkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := c.createFile(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	// write the default configuration to the file
	out, err := c.yamMarshal(DefaultConfiguration())
	if err != nil {
		return err
	}
	_, err = f.Write(out)
	return err
}
