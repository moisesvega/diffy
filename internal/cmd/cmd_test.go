package cmd

import (
	"testing"

	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/client/phabricator/phabricatormock"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/config/configmock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	oldPhabNew := phabNew
	t.Cleanup(func() {
		phabNew = oldPhabNew
	})
	phabmock := phabricatormock.NewMockClient(gomock.NewController(t))
	phabNew = func(phabricator config.Phabricator) (phabricator.Client, error) {
		return phabmock, nil
	}
	path := "diffy/settings.yaml"
	cfg := configmock.NewMockOperations(gomock.NewController(t))
	cfg.EXPECT().Read(path).Return(&config.Config{}, nil).Times(1)

	cmd, err := New(Params{
		Config: cfg,
		XDGConfigFile: func(s string) (string, error) {
			return path, nil
		},
	})
	require.NoError(t, err)
	cmd.SetArgs([]string{"--help"})
	require.NotNil(t, cmd)
	require.NotPanics(t, func() {
		err := cmd.Execute()
		require.NoError(t, err)
	})
}
