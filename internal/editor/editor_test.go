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
			return fakeExecCommand(command, args...)
		},
		In:  &in,
		Out: out,
		Err: out,
	}
	require.NoError(t, e.OpenFile(filepath))
}

const _fakeEnvVar = "GO_WANT_HELPER_PROCESS"

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	// CreateDefaults a fake command that will exit successfully
	// This wil run the TestHelperProcess function
	// and exit successfully
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	// set an environment variable to indicate that this is a fake command
	// and not a real one
	cmd.Env = []string{_fakeEnvVar + "=1"}
	return cmd
}

// TestHelperProcess isn't a real test function, it's used as a helper
// to run a fake command that will exit successfully
func TestHelperProcess(t *testing.T) {
	// This function is called by the fake command
	switch os.Getenv(_fakeEnvVar) {
	case "1":
		// This is the child process
		// called by the fake command
		os.Exit(0)
	default:
		// This is the parent process
		return
	}
}

func TestGetEditor(t *testing.T) {
	tests := []struct {
		desc string
		give string
		want string
	}{
		{
			desc: "returns vi when EDITOR not set",
			give: "",
			want: "vi",
		},
		{
			desc: "returns EDITOR value when set",
			give: "nano",
			want: "nano",
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			t.Setenv(_editorEnvVar, tt.give)
			assert.Equal(t, tt.want, getEditor())
		})
	}
}
