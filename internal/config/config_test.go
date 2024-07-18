package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadConfiguration(t *testing.T) {
	cfg, err := ReadConfiguration("testdata/config.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
}
