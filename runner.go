package main

import (
	"flag"

	"github.com/moisesvega/diffy/internal/client"
)

type runner struct {
	flagSet    *flag.FlagSet
	phabClient *client.Client
}

func (r *runner) run(args []string) error {
	r.flagSet.Parse(args)
	return nil
}
