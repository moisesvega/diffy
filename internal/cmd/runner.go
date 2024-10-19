package cmd

import (
	"os"
	"path/filepath"

	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
)

type runner struct {
	opts        opts
	editor      editor.Open
	phabricator phabricator.Client
	config      config.Operations
}

const (
	_name             = "diffy"
	_settingsFileName = "settings.yaml"
)

var (
	settingsFilePath = filepath.Join(_name, _settingsFileName)
)

func (r *runner) run(args []string) error {
	if r.opts.settings {
		return r.openAndEditConfigFile(r.editor)
	}
	return nil
}

const _XDGConfigHome = "XDG_CONFIG_HOME"

func (r *runner) openAndEditConfigFile(e editor.Open) error {
	configFilePath := filepath.Join(os.Getenv(_XDGConfigHome), settingsFilePath)
	// if the file does not exist, create it with the default configuration
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if err := r.config.CreateDefaults(configFilePath); err != nil {
			return err
		}
	}
	return e.OpenFile(configFilePath)
}
