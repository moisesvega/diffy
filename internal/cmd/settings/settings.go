package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
)

const (
	_name             = "diffy"
	_settingsFileName = "settings.yaml"
)

var settingsFilePath = filepath.Join(_name, _settingsFileName)

// NewRunner creates a new settings runner for Kong command.
func NewRunner() *runner {
	return &runner{
		xdgConfig: xdg.ConfigFile,
		editor:    editor.New(os.Stdin, os.Stdout, os.Stderr),
		config:    config.New(),
	}
}

type runner struct {
	xdgConfig func(string) (string, error)
	editor    editor.Open
	config    config.Operations
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

// Run executes the settings command for Kong.
func (r *runner) Run() error {
	sPath, err := r.xdgConfig(settingsFilePath)
	if err != nil {
		return err
	}
	return r.openAndEditConfigFile(sPath)
}
