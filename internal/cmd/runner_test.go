package cmd

import (
	"errors"
	"testing"

	"github.com/moisesvega/diffy/internal/client/github"
	"github.com/moisesvega/diffy/internal/client/github/githubmock"
	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/client/phabricator/phabricatormock"
	"github.com/moisesvega/diffy/internal/config"
	"github.com/moisesvega/diffy/internal/config/configmock"
	"github.com/moisesvega/diffy/internal/entity"
	"github.com/moisesvega/diffy/internal/entity/reportermock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRunner(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := &config.Config{}
		u := []*entity.User{
			{
				Username: "moisesvega",
				Differentials: []*entity.Differential{
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
			reporters: []entity.Reporter{rmock},
		}

		err := r.run([]string{}, "phabricator")
		require.NoError(t, err)
		require.NotNil(t, r)
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

		got := r.run([]string{}, "phabricator")
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

		got := r.run([]string{}, "phabricator")
		require.Error(t, got)
		require.ErrorIs(t, got, want)
	})

	t.Run("phab get users fail", func(t *testing.T) {
		want := errors.New("sad")
		cfg := &config.Config{}
		cmock := configmock.NewMockOperations(gomock.NewController(t))
		cmock.EXPECT().Read("").Return(cfg, nil).Times(1)

		phabmock := phabricatormock.NewMockClient(gomock.NewController(t))
		phabmock.EXPECT().GetUsers([]string{}).Return([]*entity.User{}, want).Times(1)

		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", nil
			},
			config: cmock,
			phabNew: func(phabricator config.Phabricator) (phabricator.Client, error) {
				return phabmock, nil
			},
		}

		got := r.run([]string{}, "phabricator")
		require.Error(t, got)
		require.ErrorIs(t, got, want)
	})

	t.Run("reporter fail", func(t *testing.T) {
		want := errors.New("sad")
		cfg := &config.Config{}
		cmock := configmock.NewMockOperations(gomock.NewController(t))
		cmock.EXPECT().Read("").Return(cfg, nil).Times(1)

		phabmock := phabricatormock.NewMockClient(gomock.NewController(t))
		phabmock.EXPECT().GetUsers([]string{}).Return([]*entity.User{}, nil).Times(1)
		rmock := reportermock.NewMockReporter(gomock.NewController(t))
		rmock.EXPECT().Report([]*entity.User{}).Return(want).Times(1)

		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", nil
			},
			config: cmock,
			phabNew: func(phabricator config.Phabricator) (phabricator.Client, error) {
				return phabmock, nil
			},
			reporters: []entity.Reporter{rmock},
		}

		got := r.run([]string{}, "phabricator")
		require.Error(t, got)
		require.ErrorIs(t, got, want)
	})

	t.Run("github success", func(t *testing.T) {
		cfg := &config.Config{}
		u := []*entity.User{
			{
				Username: "moisesvega",
				Differentials: []*entity.Differential{
					{
						Title:     "Fix bug",
						URI:       "https://github.com/org/repo/pull/1",
						LineCount: 42,
					},
				},
			},
		}
		cmock := configmock.NewMockOperations(gomock.NewController(t))
		cmock.EXPECT().Read("").Return(cfg, nil).Times(1)
		ghmock := githubmock.NewMockClient(gomock.NewController(t))
		ghmock.EXPECT().GetUsers([]string{"moisesvega"}).Return(u, nil).Times(1)
		rmock := reportermock.NewMockReporter(gomock.NewController(t))
		rmock.EXPECT().Report(u).Return(nil).Times(1)

		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", nil
			},
			config: cmock,
			githubNew: func(gh config.GitHub) (github.Client, error) {
				return ghmock, nil
			},
			reporters: []entity.Reporter{rmock},
		}

		err := r.run([]string{"moisesvega"}, "github")
		require.NoError(t, err)
	})

	t.Run("github new fail", func(t *testing.T) {
		want := errors.New("sad")
		cfg := &config.Config{}
		cmock := configmock.NewMockOperations(gomock.NewController(t))
		cmock.EXPECT().Read("").Return(cfg, nil).Times(1)

		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", nil
			},
			config: cmock,
			githubNew: func(gh config.GitHub) (github.Client, error) {
				return nil, want
			},
		}

		got := r.run([]string{}, "github")
		require.Error(t, got)
		require.ErrorIs(t, got, want)
	})

	t.Run("github get users fail", func(t *testing.T) {
		want := errors.New("sad")
		cfg := &config.Config{}
		cmock := configmock.NewMockOperations(gomock.NewController(t))
		cmock.EXPECT().Read("").Return(cfg, nil).Times(1)

		ghmock := githubmock.NewMockClient(gomock.NewController(t))
		ghmock.EXPECT().GetUsers([]string{}).Return([]*entity.User{}, want).Times(1)

		r := &runner{
			xdgConfig: func(s string) (string, error) {
				return "", nil
			},
			config: cmock,
			githubNew: func(gh config.GitHub) (github.Client, error) {
				return ghmock, nil
			},
		}

		got := r.run([]string{}, "github")
		require.Error(t, got)
		require.ErrorIs(t, got, want)
	})
}
