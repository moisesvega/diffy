package cmd

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFlagSet(t *testing.T) {
	give := []string{
		"--settings",
	}

	want := &opts{
		settings: true,
	}

	got := &opts{}
	pfs := pflag.NewFlagSet("new", pflag.ExitOnError)
	require.NotPanics(t, func() {
		setFlags(pfs, got)
	})
	err := pfs.Parse(give)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, want, got)
}
