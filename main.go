package main

import (
	"log"
	"os"

	"github.com/moisesvega/diffy/internal/client"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(0) // Removes timestamp

	// TODO: Make it configurable
	r := runner{cfg: &config.Config{}}
	cmd := cobra.Command{
		Use:           "diffy",
		Short:         "CLI designed to deliver comprehensive statistics and insights from code reviews and differential analysis",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return r.run(cmd.Flags().Args())
		},
	}
	registerFlags(cmd.Flags(), r.cfg)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func getUsersFromPhab(users []string) error {
	phabToken := os.Getenv("PHAB_TOKEN")
	// TODO: Make it installable to make it easier for testing
	c, err := client.New(phabToken)
	if err != nil {
		return err
	}

	_, err = c.GetUsers(users)
	if err != nil {
		return err
	}
	// report Users
	return nil
}
