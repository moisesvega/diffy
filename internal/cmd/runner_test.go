package cmd

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/client/phabricator/phabricatormock"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/config/configmock"
	"github.com/moisesvega/diffy/internal/editor/editormock"
	"github.com/moisesvega/diffy/internal/model"
	"github.com/moisesvega/diffy/internal/model/reportermock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRunner(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := &config.Config{}
		u := []*model.User{
			{
				Username: "moisesvega",
				Differentials: []*model.Differential{
					{
						Title:     "title",
						URI:       "uri",
						LineCount: 11,
					},
				},
			},
		}
		cmock := configmock.NewMockOperations(gomock.NewController(t))
		cmock.EXPECT().Read("").Return(cfg, nil).Times(1)
		phabmock := phabricatormock.NewMockClient(gomock.NewController(t))
		phabmock.EXPECT().GetUsers([]string{}).Return(u, nil).Times(1)
		rmock := reportermock.NewMockReporter(gomock.NewController(t))
		rmock.EXPECT().Report(u).Return(nil).Times(1)

		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", nil
			},
			config: cmock,
			phabNew: func(phabricator config.Phabricator) (phabricator.Client, error) {
				return phabmock, nil
			},
			reporters: []model.Reporter{rmock},
		}

		err := r.run([]string{})
		require.NoError(t, err)
		require.NotNil(t, r)
	})
	t.Run("xdgConfig error", func(t *testing.T) {
		want := errors.New("sad")
		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", want
			},
		}
		err := r.run([]string{})
		require.Error(t, err)
		assert.ErrorIs(t, err, want)
	})

	t.Run("config fail", func(t *testing.T) {
		want := errors.New("sad")
		cmock := configmock.NewMockOperations(gomock.NewController(t))
		cmock.EXPECT().Read("").Return(&config.Config{}, want).Times(1)

		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", nil
			},
			config: cmock,
		}

		got := r.run([]string{})
		require.Error(t, got)
		require.ErrorIs(t, got, want)
	})

	t.Run("phab new fail", func(t *testing.T) {
		want := errors.New("sad")
		cfg := &config.Config{}
		cmock := configmock.NewMockOperations(gomock.NewController(t))
		cmock.EXPECT().Read("").Return(cfg, nil).Times(1)

		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", nil
			},
			config: cmock,
			phabNew: func(phabricator config.Phabricator) (phabricator.Client, error) {
				return nil, want
			},
		}

		got := r.run([]string{})
		require.Error(t, got)
		require.ErrorIs(t, got, want)
	})

	t.Run("phab get users fail", func(t *testing.T) {
		want := errors.New("sad")
		cfg := &config.Config{}
		cmock := configmock.NewMockOperations(gomock.NewController(t))
		cmock.EXPECT().Read("").Return(cfg, nil).Times(1)

		phabmock := phabricatormock.NewMockClient(gomock.NewController(t))
		phabmock.EXPECT().GetUsers([]string{}).Return([]*model.User{}, want).Times(1)

		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", nil
			},
			config: cmock,
			phabNew: func(phabricator config.Phabricator) (phabricator.Client, error) {
				return phabmock, nil
			},
		}

		got := r.run([]string{})
		require.Error(t, got)
		require.ErrorIs(t, got, want)
	})

	t.Run("reporter fail", func(t *testing.T) {
		want := errors.New("sad")
		cfg := &config.Config{}
		cmock := configmock.NewMockOperations(gomock.NewController(t))
		cmock.EXPECT().Read("").Return(cfg, nil).Times(1)

		phabmock := phabricatormock.NewMockClient(gomock.NewController(t))
		phabmock.EXPECT().GetUsers([]string{}).Return([]*model.User{}, nil).Times(1)
		rmock := reportermock.NewMockReporter(gomock.NewController(t))
		rmock.EXPECT().Report([]*model.User{}).Return(want).Times(1)

		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", nil
			},
			config: cmock,
			phabNew: func(phabricator config.Phabricator) (phabricator.Client, error) {
				return phabmock, nil
			},
			reporters: []model.Reporter{rmock},
		}

		got := r.run([]string{})
		require.Error(t, got)
		require.ErrorIs(t, got, want)
	})
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
