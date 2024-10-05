package main

import (
	"log"

	"github.com/moisesvega/diffy/internal/cmd"
)

func main() {
	if err := cmd.Main().Execute(); err != nil {
		log.Fatal(err)
	}
}
