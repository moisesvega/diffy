package cmd

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
)

type runner struct {
	opts        opts
	phabricator phabricator.Client
	editor      editor.Open
}

const (
	_name             = "diffy"
	_settingsFileName = "settings.yaml"
)

var (
	settingsFilePath = filepath.Join(_name, _settingsFileName)
)

func (r *runner) run(args []string, cfg *config.Config) error {
	if r.opts.Settings {
		return openAndEditConfigFile()
	}
	return nil
}

func openAndEditConfigFile() error {
	configFilePath, err := xdg.ConfigFile(settingsFilePath)
	if err != nil {
		return err
	}
	return editor.New(os.Stdin, os.Stdout, os.Stderr).OpenFile(configFilePath)
}
