package cmd

import (
	"io"
	"log"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
)

type runner struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	cfg         *config.Config
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

func (r *runner) run(args []string) (err error) {
	if r.cfg.Settings {
		configFilePath, err := xdg.ConfigFile(settingsFilePath)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Save the config file at:", configFilePath)
		return r.editor.OpenFile(configFilePath)
	}
	return nil
}
