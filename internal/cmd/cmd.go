package cmd

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/alecthomas/kong"
	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/cmd/settings"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
	"github.com/moisesvega/diffy/internal/entity"
	"github.com/moisesvega/diffy/internal/reporter/heatmap"
)

type CLI struct {
	Analyze  AnalyzeCmd  `cmd:"analyze" help:"Analyze users and generate reports."`
	Settings SettingsCmd `cmd:"settings" help:"Opens the settings file in the editor."`
}

type AnalyzeCmd struct {
	Users []string `arg:"" optional:"" help:"Usernames to analyze (e.g., user1 user2)."`
}

func (c *AnalyzeCmd) Run() error {
	r := runner{
		editor:    editor.New(os.Stdin, os.Stdout, os.Stderr),
		config:    config.New(),
		phabNew:   phabricator.New,
		xdgConfig: xdg.ConfigFile,
		reporters: []entity.Reporter{heatmap.New()},
	}
	return r.run(c.Users)
}

type SettingsCmd struct{}

func (c *SettingsCmd) Run() error {
	r := settings.NewRunner()
	return r.Run()
}

func Main() *kong.Kong {
	cli := &CLI{}
	return kong.Must(cli,
		kong.Name("diffy"),
		kong.Description("CLI designed to deliver comprehensive statistics and insights from code reviews and differential analysis"),
		kong.UsageOnError(),
	)
}
