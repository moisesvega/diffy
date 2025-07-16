package config

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v3"
)

func TestReadConfiguration(t *testing.T) {
	file := `
apis:
    phabricator:
        base_url: <replace_me>
        api_token: ""
        api_token_env: PHAB_API_TOKEN
        access_token: ""
        access_token_env: PHAB_ACCESS_TOKEN
        arrc_file_path: ~/.arcrc
    github:
        base_url: https://github.com/
        api_token: ""
        api_token_env: GITHUB_API_TOKEN
`
	fp := path.Join(t.TempDir(), "test.yaml")
	require.NoError(t, os.WriteFile(fp, []byte(file), _mode))
	ctrl := New()
	cfg, err := ctrl.Read(fp)
	require.NoError(t, err)
	require.NotEmpty(t, cfg)
}

func TestReadConfigurationParseError(t *testing.T) {
	file := `
phab_users:
  - phabricator_config: {json: - fail}
`
	fp := path.Join(t.TempDir(), "test.yaml")
	require.NoError(t, os.WriteFile(fp, []byte(file), _mode))
	ctrl := New()
	cfg, err := ctrl.Read(fp)
	require.Error(t, err)
	require.Nil(t, cfg)
}

func TestReadConfigurationFileNotFound(t *testing.T) {
	ctrl := New()
	fp := path.Join(t.TempDir(), "test.yaml")
	cfg, err := ctrl.Read(fp)
	require.Error(t, err)
	require.Nil(t, cfg)
}

func TestWriteConfiguration(t *testing.T) {
	defaults := DefaultConfiguration()
	got, err := yaml.Marshal(defaults)
	require.NoError(t, err)
	require.NotEmpty(t, got)
	// TODO(moisesvega): Uncomment this line to update the defaults.yaml file
	// CreateDefaults a small script to run this test and update the defaults.yaml file
	// os.WriteFile(path.Join("./testdata/defaults.yaml"), got, _mode)
	want, err := os.ReadFile("./testdata/defaults.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, want)
	require.Equal(t, string(want), string(got))
}

func TestController_ReadConfiguration(t *testing.T) {
	dir := t.TempDir()
	t.Run("success", func(t *testing.T) {
		ctrl := New()
		fp := path.Join(dir, "test.yaml")
		require.NoError(t, ctrl.CreateDefaults(fp))
		assert.FileExists(t, fp)
		got, err := ctrl.Read(fp)
		require.NoError(t, err)
		want := DefaultConfiguration()
		assert.EqualValues(t, want, got)
	})

	t.Run("fail while creating directory", func(t *testing.T) {
		ctrl := New()
		want := errors.New("sad")
		//  mock mkdirAll to return an error
		fp := path.Join(dir, "not_found.yaml")
		ctrl.mkdirAll = func(path string, perm os.FileMode) error {
			assert.Equal(t, dir, path)
			assert.Equal(t, fs.FileMode(0o755), perm)
			return want
		}
		err := ctrl.CreateDefaults(fp)
		require.Error(t, err)
	})

	t.Run("fail while creating file", func(t *testing.T) {
		ctrl := New()
		want := errors.New("sad")
		//  mock mkdirAll to return an error
		fp := path.Join(dir, "not_found.yaml")
		ctrl.mkdirAll = func(path string, perm os.FileMode) error {
			return nil
		}
		ctrl.createFile = func(name string) (*os.File, error) {
			assert.Equal(t, fp, name)
			return nil, want
		}
		err := ctrl.CreateDefaults(fp)
		require.Error(t, err)
	})

	t.Run("fail while marshaling the config", func(t *testing.T) {
		ctrl := New()
		want := errors.New("sad")
		//  mock mkdirAll to return an error
		fp := path.Join(dir, "not_found.yaml")
		ctrl.mkdirAll = func(path string, perm os.FileMode) error {
			return nil
		}
		ctrl.createFile = func(name string) (*os.File, error) {
			assert.Equal(t, fp, name)
			return &os.File{}, nil
		}
		ctrl.yamMarshal = func(in interface{}) ([]byte, error) {
			return nil, want
		}
		err := ctrl.CreateDefaults(fp)
		require.Error(t, err)
	})
}
