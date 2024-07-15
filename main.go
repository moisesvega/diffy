package main

import (
	"log"
	"os"

	"github.com/moisesvega/diffy/internal/client"
)

func main() {
	log.SetFlags(0) // Removes timestamp

	// TODO: Make it configurable
	phabToken := os.Getenv("PHAB_TOKEN")
	if err := run(phabToken, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(phabToken string, users []string) error {
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
