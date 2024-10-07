package cmd

import (
	"io"
	"log"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/phabricator"
)

type runner struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	cfg         *config.Config
	phabricator phabricator.Client
}

const (
	_name             = "diffy"
	_settingsFileName = "settings.yaml"
)

var (
	_settingsFilePath = filepath.Join(_name, _settingsFileName)
)

func (r *runner) run(args []string) (err error) {
	if r.cfg.Settings {
		configFilePath, err := xdg.ConfigFile(_settingsFilePath)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Save the config file at:", configFilePath)

		return openSettings(r.stdin, r.stdout, r.stderr, configFilePath)
	}
	if r.phabricator == nil {
		r.phabricator, err = phabricator.New(&r.cfg.APIs.Phabricator)
		if err != nil {
			return err
		}
	}
	return nil
}
