package cmd

import (
	"errors"
	"testing"

	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/config/configmock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRunner(t *testing.T) {
	t.Run("xdg fail", func(t *testing.T) {
		want := errors.New("sad")
		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", want
			},
		}

		_, got := r.getPhabricatorClient()
		require.Error(t, got)
		require.ErrorIs(t, got, want)
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

		_, got := r.getPhabricatorClient()
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

		_, got := r.getPhabricatorClient()
		require.Error(t, got)
		require.ErrorIs(t, got, want)
	})
}

func TestRunE(t *testing.T) {
	r := &runner{}
	require.NoError(t, r.runE(&cobra.Command{}, []string{}))
}
