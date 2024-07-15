package main

import (
	"testing"

	"github.com/moisesvega/diffy/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFlagSet(t *testing.T) {
	tests := []struct {
		desc string
		give []string
		want *config.Config
	}{
		{desc: "no flags", give: []string{}, want: &config.Config{}},
		{desc: "settings", give: []string{"--settings"}, want: &config.Config{}},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cfg := &config.Config{}
			fs := newFlagSet("test", cfg, nil)
			err := fs.Parse(tt.give)
			require.NoError(t, err)
			require.NotNil(t, cfg)
			assert.Equal(t, tt.want, cfg)
		})
	}
}
