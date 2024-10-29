package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
	"github.com/spf13/cobra"
)

const (
	_name             = "diffy"
	_settingsFileName = "settings.yaml"
)

var (
	settingsFilePath = filepath.Join(_name, _settingsFileName)
)

// New returns a command that opens the settings file in the editor.
func New() *cobra.Command {
	r := runner{
		xdgConfig: xdg.ConfigFile,
		editor:    editor.New(os.Stdin, os.Stdout, os.Stderr),
		config:    config.New(),
	}
	cmd := &cobra.Command{
		Use:           "settings",
		Short:         "Opens the settings file in the editor. By default, uses $XDG_CONFIG_HOME as the path. On macOS, if $XDG_CONFIG_HOME is not set, defaults to $HOME/Library/Application Support/diffy",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          r.runE,
	}
	return cmd
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

func (r *runner) runE(cmd *cobra.Command, args []string) error {
	sPath, err := r.xdgConfig(settingsFilePath)
	if err != nil {
		return err
	}
	return r.openAndEditConfigFile(sPath)
}
