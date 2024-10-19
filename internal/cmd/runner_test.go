package cmd

import (
	"path/filepath"
	"testing"

	"github.com/moisesvega/diffy/internal/config/configmock"
	"github.com/moisesvega/diffy/internal/editor/editormock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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
	emock := editormock.NewMockOpen(gomock.NewController(t))

	pathToSettings := filepath.Join(tmpDir, settingsFilePath)
	emock.EXPECT().OpenFile(pathToSettings).Return(nil).Times(1)

	cfgmock := configmock.NewMockOperations(gomock.NewController(t))
	cfgmock.EXPECT().CreateDefaults(pathToSettings).Return(nil).Times(1)
	r := &runner{
		opts: opts{
			settings: true,
		},
		editor: emock,
		config: cfgmock,
	}

	err := r.run([]string{})
	require.NoError(t, err)
	require.NotNil(t, r)
}
