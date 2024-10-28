package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
	"github.com/moisesvega/diffy/internal/filter"
	"github.com/moisesvega/diffy/internal/model"
)

type runner struct {
	opts      opts
	editor    editor.Open
	phabNew   func(config.Phabricator) (phabricator.Client, error)
	config    config.Operations
	xdgConfig func(string) (string, error)
	reporters []model.Reporter
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

	cfg, err := r.config.Read(sPath)
	if err != nil {
		return fmt.Errorf("%w: %w", errConfigNotFound, err)
	}

	phab, err := r.phabNew(cfg.APIs.Phabricator)
	if err != nil {
		return err
	}

	u, err := phab.GetUsers(args)
	if err != nil {
		return err
	}

	// filter out closed differentials and those with less than 10 lines
	for _, user := range u {
		diffs := slices.DeleteFunc(user.Differentials, filter.ByStatus(model.Closed))
		user.Differentials = slices.DeleteFunc(diffs, filter.MinLineCount(10))
	}

	// report the data
	for _, reporter := range r.reporters {
		if err := reporter.Report(u); err != nil {
			return fmt.Errorf("failed to report: %w", err)
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
