package settings

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/moisesvega/diffy/internal/config/configmock"
	"github.com/moisesvega/diffy/internal/editor/editormock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewRunner(t *testing.T) {
	require.NotNil(t, NewRunner())
}

func TestRun(t *testing.T) {
	r := &runner{
		xdgConfig: func(s string) (string, error) {
			return "", nil
		},
	}
	require.NotNil(t, r)
}

func TestXDGConfigError(t *testing.T) {
	want := errors.New("sad")
	r := &runner{
		xdgConfig: func(s string) (string, error) {
			return "", want
		},
	}
	err := r.Run()
	require.Error(t, err)
	assert.ErrorIs(t, err, want)
}

func TestOpenSettings(t *testing.T) {
	tmpDir := t.TempDir()
	pathToSettings := filepath.Join(tmpDir, settingsFilePath)
	emock := editormock.NewMockOpen(gomock.NewController(t))
	emock.EXPECT().OpenFile(pathToSettings).Return(nil).Times(1)
	cfgmock := configmock.NewMockOperations(gomock.NewController(t))
	cfgmock.EXPECT().CreateDefaults(pathToSettings).Return(nil).Times(1)
	r := &runner{
		editor: emock,
		config: cfgmock,
		xdgConfig: func(s string) (string, error) {
			return pathToSettings, nil
		},
	}

	err := r.Run()
	require.NoError(t, err)
	require.NotNil(t, r)
}

func TestOpenSettingsError(t *testing.T) {
	tmpDir := t.TempDir()
	want := errors.New("sad")
	pathToSettings := filepath.Join(tmpDir, settingsFilePath)
	emock := editormock.NewMockOpen(gomock.NewController(t))
	cfgmock := configmock.NewMockOperations(gomock.NewController(t))
	cfgmock.EXPECT().CreateDefaults(pathToSettings).Return(want).Times(1)
	r := &runner{
		editor: emock,
		config: cfgmock,
		xdgConfig: func(s string) (string, error) {
			return pathToSettings, nil
		},
	}
	err := r.Run()
	require.Error(t, err)
	assert.ErrorIs(t, err, want)
}
