package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCMD(t *testing.T) {
	require.NotPanics(t, func() {
		oldArgs := os.Args
		t.Cleanup(func() {
			os.Args = oldArgs
		})
		os.Args = []string{"cmd", "--help"}
		main()
	})
}
