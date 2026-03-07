package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/moisesvega/diffy/internal/client/github"
	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
	"github.com/moisesvega/diffy/internal/entity"
	"github.com/moisesvega/diffy/internal/filter"
)

type runner struct {
	editor    editor.Open
	phabNew   func(config.Phabricator) (phabricator.Client, error)
	githubNew func(config.GitHub) (github.Client, error)
	config    config.Operations
	xdgConfig func(string) (string, error)
	reporters []entity.Reporter
}

const (
	_name             = "diffy"
	_settingsFileName = "settings.yaml"
)

var (
	settingsFilePath  = filepath.Join(_name, _settingsFileName)
	errConfigNotFound = errors.New("config file not found, please run `diffy --settings` to create one")
)

type userFetcher interface {
	GetUsers(names []string) ([]*entity.User, error)
}

func (r *runner) run(args []string, source string) error {
	sPath, err := r.xdgConfig(settingsFilePath)
	if err != nil {
		return err
	}
	cfg, err := r.config.Read(sPath)
	if err != nil {
		return fmt.Errorf("%w: %w", errConfigNotFound, err)
	}

	var fetcher userFetcher
	switch source {
	case "github":
		fetcher, err = r.githubNew(cfg.APIs.GitHub)
	default:
		fetcher, err = r.phabNew(cfg.APIs.Phabricator)
	}
	if err != nil {
		return err
	}

	u, err := fetcher.GetUsers(args)
	if err != nil {
		return err
	}

	// filter out closed differentials and those with less than 10 lines
	for _, user := range u {
		diffs := slices.DeleteFunc(user.Differentials, filter.ByStatus(entity.Closed))
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
