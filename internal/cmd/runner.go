package cmd

import (
	"errors"
	"fmt"
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
	xdgConfig   func(string) (string, error)
	cfg         *config.Config
}

const (
	_name             = "diffy"
	_settingsFileName = "settings.yaml"
)

var (
	settingsFilePath  = filepath.Join(_name, _settingsFileName)
	errConfigNotFound = errors.New("config file not found, please run `diffy --settings` to create one")
)

func (r *runner) run(args []string) error {
	sPath, err := r.xdgConfig(settingsFilePath)
	if err != nil {
		return err
	}
	if r.opts.settings {
		return r.openAndEditConfigFile(sPath)
	}

	u, err := r.phabricator.GetUsers([]string{r.cfg.Me.PhabricatorUsername})
	for _, user := range u {
		for _, differential := range user.Differentials {
			fmt.Println(differential)
		}
	}
	return nil
}

func (r *runner) openAndEditConfigFile(path string) error {
	// if the file does not exist, create it with the default configuration
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Creating a new configuration file at", settingsFilePath)
		if err := r.config.CreateDefaults(path); err != nil {
			return err
		}
	}

	return r.editor.OpenFile(path)
}
