package cmd

import (
	"testing"

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
