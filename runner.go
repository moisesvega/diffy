package main

import (
	"log"
	"strings"

	"github.com/moisesvega/diffy/internal/client"
	"github.com/moisesvega/diffy/internal/config"
)

type runner struct {
	cfg        *config.Config
	phabClient *client.Client
}

func (r *runner) run(args []string) error {
	if r.cfg.Settings {
		return editSettings()
	}
	log.Printf("Running %s", strings.Join(args, " "))
	log.Println(r.cfg)
	return nil
}

func editSettings() error {
	// TODO: Open settings.yaml and a allow user to edit it.
	return nil
}
