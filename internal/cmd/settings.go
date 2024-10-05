package cmd

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func openSettings(in io.Reader, out, err io.Writer, filepath string) error {
	v := os.Getenv("EDITOR")
	// TODO(moises): make sure
	log.Println(filepath)
	cmd := exec.Command(v, filepath)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = in, out, err
	return cmd.Run()
}
