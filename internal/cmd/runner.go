package cmd

import (
	"os"
	"path/filepath"

	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
	"gopkg.in/yaml.v3"
)

type runner struct {
	opts        opts
	editor      editor.Open
	phabricator phabricator.Client
}

const (
	_name             = "diffy"
	_settingsFileName = "settings.yaml"
)

var (
	settingsFilePath = filepath.Join(_name, _settingsFileName)
)

func (r *runner) run(args []string, cfg *config.Config) error {
	if r.opts.settings {
		return openAndEditConfigFile(r.editor)
	}
	return nil
}

const _XDGConfigHome = "XDG_CONFIG_HOME"

func openAndEditConfigFile(e editor.Open) error {
	configFilePath := filepath.Join(os.Getenv(_XDGConfigHome), settingsFilePath)
	// if the file does not exist, create it with the default configuration
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if err := createDefaultConfigFile(configFilePath); err != nil {
			return err
		}
	}
	return e.OpenFile(configFilePath)
}

func createDefaultConfigFile(path string) error {
	// create the file
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	// write the default configuration
	config.DefaultConfiguration()
	out, err := yaml.Marshal(config.DefaultConfiguration())
	if err != nil {
		return err
	}
	_, err = f.Write(out)
	return err
}
