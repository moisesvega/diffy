package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCLI(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() {
		os.Args = oldArgs
	})
	os.Args = []string{
		"-h",
	}
	require.NotPanics(t, func() {
		main()
	})
}
