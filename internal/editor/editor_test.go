package editor

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	var in bytes.Buffer
	out := &bytes.Buffer{}
	e := New(&in, out, out)
	require.NotNil(t, e)
	require.Equal(t, &in, e.In)
	require.Equal(t, out, e.Out)
	require.Equal(t, out, e.Err)
}

func TestOpenFile(t *testing.T) {
	editor := "nvim"
	filepath := "test.txt"
	t.Setenv(_editorEnvVar, editor)
	var in bytes.Buffer
	out := &bytes.Buffer{}
	e := &Controller{
		command: func(command string, args ...string) *exec.Cmd {
			assert.Equal(t, editor, command)
			require.Len(t, args, 1)
			assert.Equal(t, filepath, args[0])
			// Return a fake command that will exit successfully
			return fakeExecCommand("echo", "success")
		},
		In:  &in,
		Out: out,
		Err: out,
	}
	require.NoError(t, e.OpenFile(filepath))
}

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}
