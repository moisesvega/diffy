package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/spf13/cobra"
)

type runner struct {
	phabNew   func(config.Phabricator) (phabricator.Client, error)
	config    config.Operations
	xdgConfig func(string) (string, error)
}

const (
	_name             = "diffy"
	_settingsFileName = "settings.yaml"
)

var (
	settingsFilePath  = filepath.Join(_name, _settingsFileName)
	errConfigNotFound = errors.New("config file not found, please run `diffy --settings` to create one")
)

func (r *runner) runE(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func (r *runner) getPhabricatorClient() (phabricator.Client, error) {
	sPath, err := r.xdgConfig(settingsFilePath)
	if err != nil {
		return nil, err
	}
	cfg, err := r.config.Read(sPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errConfigNotFound, err)
	}
	return r.phabNew(cfg.APIs.Phabricator)
}
