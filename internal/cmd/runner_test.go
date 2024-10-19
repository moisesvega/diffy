package cmd

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

func TestRunner(t *testing.T) {
	want := errors.New("sad")
	r := &runner{
		xdgConfig: func(s string) (string, error) {
			return "", want
		},
	}
	err := r.run([]string{})
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
		opts: opts{
			settings: true,
		},
		editor: emock,
		config: cfgmock,
		xdgConfig: func(s string) (string, error) {
			return pathToSettings, nil
		},
	}

	err := r.run([]string{})
	require.NoError(t, err)
	require.NotNil(t, r)
}

func TestOpenSettingsError(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tmpDir := t.TempDir()

		pathToSettings := filepath.Join(tmpDir, settingsFilePath)
		emock := editormock.NewMockOpen(gomock.NewController(t))

		want := errors.New("sad")
		cfgmock := configmock.NewMockOperations(gomock.NewController(t))
		cfgmock.EXPECT().CreateDefaults(pathToSettings).Return(want).Times(1)
		r := &runner{
			opts: opts{
				settings: true,
			},
			editor: emock,
			config: cfgmock,
			xdgConfig: func(s string) (string, error) {
				return pathToSettings, nil
			},
		}

		err := r.run([]string{})
		require.Error(t, err)
		assert.ErrorIs(t, err, want)
	})
}
