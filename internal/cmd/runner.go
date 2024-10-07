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

	cfg        *config.Config
	phabClient phabricator.Client
}

const (
	_name             = "diffy"
	_settingsFileName = "settings.yaml"
)

func (r *runner) run(args []string) error {
	if r.cfg.Settings {
		// TODO(moises): make it configurable.
		// Finding application config files.
		// SearchConfigFile takes one parameter which must contain the name of
		// the file, but it can also contain a set of parent directories relative
		// to the config search paths (xdg.ConfigHome and xdg.ConfigDirs).
		configFilePath, err := xdg.SearchConfigFile("appname/config.yaml")
		if err != nil {
			// TODO(moises) create a default settings.yaml
			log.Fatal(err)
		}

		configFilePath, err = xdg.ConfigFile(filepath.Join(_name, _settingsFileName))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Save the config file at:", configFilePath)

		return openSettings(r.stdin, r.stdout, r.stderr, configFilePath)
	}
	return nil
}
