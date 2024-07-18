package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestReadConfiguration(t *testing.T) {
	file := `
phab_users:
  - user
  - user1
github_users:
  - user
  - user1
phabricator_config:
  url: "someURL"
`
	fp := path.Join(t.TempDir(), "test.yaml")
	require.NoError(t, os.WriteFile(fp, []byte(file), 0644))
	cfg, err := ReadConfiguration(fp)
	require.NoError(t, err)
	require.NotEmpty(t, cfg)
}

func TestReadConfigurationParseError(t *testing.T) {
	file := `
phab_users:
  - phabricator_config:
`
	fp := path.Join(t.TempDir(), "test.yaml")
	require.NoError(t, os.WriteFile(fp, []byte(file), 0644))
	cfg, err := ReadConfiguration(fp)
	require.Error(t, err)
	require.Nil(t, cfg)
}

func TestReadConfigurationFileNotFound(t *testing.T) {
	fp := path.Join(t.TempDir(), "test.yaml")
	cfg, err := ReadConfiguration(fp)
	require.Error(t, err)
	require.Nil(t, cfg)
}

func TestWriteConfiguration(t *testing.T) {
	defaults := DefaultConfiguration()
	got, err := yaml.Marshal(defaults)
	require.NoError(t, err)
	require.NotEmpty(t, got)
	// os.WriteFile(path.Join("./testdata/defaults.yaml"), got, 0644)
	want, err := os.ReadFile("./testdata/defaults.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, want)
	require.Equal(t, string(got), string(want))
}
