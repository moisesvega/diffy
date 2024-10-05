package cmd

import (
	"github.com/moisesvega/diffy/internal/config"
	"github.com/spf13/cobra"
)

func Main() *cobra.Command {
	r := &runner{
		cfg: &config.Config{},
	}
	cmd := &cobra.Command{
		Use:           "diffy",
		Short:         "CLI designed to deliver comprehensive statistics and insights from code reviews and differential analysis",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return r.run(cmd.Flags().Args())
		},
	}
	setFlags(cmd.Flags(), r.cfg)
	return cmd
}
