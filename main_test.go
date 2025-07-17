package main

import (
	"testing"

	"github.com/moisesvega/diffy/internal/cmd"
	"github.com/stretchr/testify/require"
)

func TestNewCMD(t *testing.T) {
	// Test that the command structure initializes correctly
	require.NotPanics(t, func() {
		// This tests that cmd.Main() doesn't panic during construction
		kong := cmd.Main(_version)
		require.NotNil(t, kong)
	})
}
