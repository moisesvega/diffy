package cmd

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor"
	"github.com/spf13/cobra"
)

func Main() *cobra.Command {
	o := &opts{}
	cmd := &cobra.Command{
		Use:           "diffy",
		Short:         "CLI designed to deliver comprehensive statistics and insights from code reviews and differential analysis",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			r := runner{
				opts:      *o,
				editor:    editor.New(os.Stdin, os.Stdout, os.Stderr),
				config:    config.New(),
				phabNew:   phabricator.New,
				xdgConfig: xdg.ConfigFile,
			}
			return r.run(cmd.Flags().Args())
		},
	}
	setFlags(cmd.Flags(), o)
	return cmd
}
