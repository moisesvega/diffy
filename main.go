package main

import (
	"log"

	"github.com/adrg/xdg"
	"github.com/moisesvega/diffy/internal/cmd"
	"github.com/moisesvega/diffy/internal/config"
)

func main() {
	c, err := cmd.New(cmd.Params{
		Config:        config.New(),
		XDGConfigFile: xdg.ConfigFile,
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
