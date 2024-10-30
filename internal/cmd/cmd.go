package cmd

import (
	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/cmd/contributions"
	"github.com/moisesvega/diffy/internal/cmd/settings"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/spf13/cobra"
)

type Params struct {
	Config        config.Operations
	XDGConfigFile func(string) (string, error)
}

var phabNew = phabricator.New

func New(p Params) (*cobra.Command, error) {
	r := runner{
		config:    p.Config,
		phabNew:   phabNew,
		xdgConfig: p.XDGConfigFile,
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
	c, _ := r.getPhabricatorClient()
	cmd.AddCommand(contributions.New(c))
	return cmd, nil
}
