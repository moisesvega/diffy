package editor

import (
	"io"
	"os"
	"os/exec"
)

type Open interface {
	OpenFile(string) error
}

// Controller manages actions related to the editor.
type Controller struct {
	In       io.Reader
	Out, Err io.Writer

	// command is the function used to run the editor command
	// It is a variable in order to be able to mock it in tests
	command func(command string, args ...string) *exec.Cmd
}

// New creates a new Controller.
func New(in io.Reader, out, err io.Writer) *Controller {
	return &Controller{
		In:      in,
		Out:     out,
		Err:     err,
		command: exec.Command,
	}
}

const _editorEnvVar = "EDITOR"

// OpenFile opens a file in the editor
// It uses the EDITOR environment variable to determine the editor to use
func (e *Controller) OpenFile(filepath string) error {
	cmd := e.command(os.Getenv(_editorEnvVar), filepath)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = e.In, e.Out, e.Err
	return cmd.Run()
}
