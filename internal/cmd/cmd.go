package cmd

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/cmd/settings"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
	"github.com/moisesvega/diffy/internal/entity"
	"github.com/moisesvega/diffy/internal/reporter/heatmap"
	"github.com/spf13/cobra"
)

func Main() *cobra.Command {
	r := runner{
		editor:    editor.New(os.Stdin, os.Stdout, os.Stderr),
		config:    config.New(),
		phabNew:   phabricator.New,
		xdgConfig: xdg.ConfigFile,
		reporters: []entity.Reporter{heatmap.New()},
	}
	cmd := &cobra.Command{
		Use:           "diffy",
		Short:         "CLI designed to deliver comprehensive statistics and insights from code reviews and differential analysis",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          r.runE,
	}
	cmd.AddCommand(settings.New())
	return cmd
}

func (r *runner) runE(cmd *cobra.Command, args []string) error {
	return r.run(cmd.Flags().Args())
}
