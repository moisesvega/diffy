package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
	"github.com/moisesvega/diffy/internal/model"
	"github.com/moisesvega/diffy/internal/reporter/heatmap"
	"github.com/uber/gonduit/constants"
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

	phab, err := phabricator.New(r.cfg.APIs.Phabricator)
	if err != nil {
		return err
	}

	u, err := phab.GetUsers(args)
	if err != nil {
		return err
	}

	// TODO(moisesvega): Create filters instead
	for _, user := range u {
		closed := make([]*model.Differential, 0)
		for _, differential := range user.Differentials {
			if differential.Status == constants.DifferentialStatusLegacyPublished {
				closed = append(closed, differential)
			}
		}
		user.Differentials = closed
	}
	// TODO(moisesvega): Make it configurable
	hm := heatmap.New()
	return hm.Report(u)
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
