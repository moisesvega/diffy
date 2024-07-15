package main

import (
	"flag"
	"io"

	"github.com/moisesvega/diffy/internal/config"
)

func newFlagSet(name string, cfg *config.Config, output io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(output)
	fs.Bool("settings", cfg.Settings, "edit configurations")
	return fs
}
