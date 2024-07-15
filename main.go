package main

import (
	"flag"
	"log"
	"os"

	"github.com/moisesvega/diffy/internal/client"
)

func main() {
	log.SetFlags(0) // Removes timestamp

	// TODO: Make it configurable
	r := runner{
		flagSet: flag.NewFlagSet("diffy", flag.ContinueOnError),
	}
	if err := r.run(os.Args[1:]); err != nil {
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
