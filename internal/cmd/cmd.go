package cmd

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
	"github.com/moisesvega/diffy/internal/model"
	"github.com/moisesvega/diffy/internal/reporter/heatmap"
	"github.com/spf13/cobra"
)

func Main() *cobra.Command {
	r := runner{
		editor:    editor.New(os.Stdin, os.Stdout, os.Stderr),
		config:    config.New(),
		phabNew:   phabricator.New,
		xdgConfig: xdg.ConfigFile,
		reporters: []model.Reporter{heatmap.New()},
	}
	cmd := &cobra.Command{
		Use:           "diffy",
		Short:         "CLI designed to deliver comprehensive statistics and insights from code reviews and differential analysis",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	o := &opts{}
	setFlags(cmd.Flags(), o)
	cmd.RunE = runE(&r, o)
	return cmd
}

func runE(r *runner, o *opts) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		r.opts = *o
		return r.run(cmd.Flags().Args())
	}
}
