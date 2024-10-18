package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/editor/editormock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gopkg.in/yaml.v3"
)

func TestRunner(t *testing.T) {
	r := &runner{}
	err := r.run([]string{})
	require.NoError(t, err)
}

func TestOpenSettings(t *testing.T) {

	// Set the XDG_CONFIG_HOME environment variable to a temporary directory
	// xdg.ConfigFile will use this directory to create the settings file
	tmpDir := t.TempDir()
	t.Setenv(_XDGConfigHome, tmpDir)
	mock := editormock.NewMockOpen(gomock.NewController(t))

	pathToSettings := filepath.Join(tmpDir, settingsFilePath)
	mock.EXPECT().OpenFile(pathToSettings).Return(nil).Times(1)

	r := &runner{
		opts: opts{
			settings: true,
		},
		editor: mock,
		config: config.New(),
	}

	err := r.run([]string{})
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.FileExists(t, pathToSettings)
	f, err := os.ReadFile(pathToSettings)
	require.NoError(t, err)
	// TODO(moisesvega): Update the test and use fake configuration controller
	got := &config.Config{}
	require.NoError(t, yaml.Unmarshal(f, got))
	assert.EqualValues(t, config.DefaultConfiguration(), got)
}
