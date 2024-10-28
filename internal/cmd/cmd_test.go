package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCMD(t *testing.T) {
	cmd := Main()
	cmd.SetArgs([]string{"--help"})
	require.NotNil(t, cmd)
	require.NotPanics(t, func() {
		err := cmd.Execute()
		require.NoError(t, err)
	})
}

func TestRunE(t *testing.T) {
	want := errors.New("sad")
	r := &runner{
		xdgConfig: func(s string) (string, error) {
			return "", want
		},
	}
	f := runE(r, &opts{})
	require.NotNil(t, f)
	require.NotPanics(t, func() {
		err := f(&cobra.Command{}, []string{})
		require.Error(t, err)
		assert.ErrorIs(t, err, want)
	})
}
