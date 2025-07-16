package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCMD(t *testing.T) {
	kong := Main()
	require.NotNil(t, kong)
}

func TestRun(t *testing.T) {
	want := errors.New("sad")
	r := &runner{
		xdgConfig: func(s string) (string, error) {
			return "", want
		},
	}
	require.NotPanics(t, func() {
		err := r.run([]string{})
		require.Error(t, err)
		assert.ErrorIs(t, err, want)
	})
}
