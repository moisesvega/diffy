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
me:
    phabricator: <replace_me>
    github: <replace_me>
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
teams:
    a_team:
        phabricator_users:
            - <replace_me>
        github_users:
            - <replace_me>
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
  - phabricator_config: {json: - fail}
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
	os.WriteFile(path.Join("./testdata/defaults.yaml"), got, 0644)
	want, err := os.ReadFile("./testdata/defaults.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, want)
	require.Equal(t, string(got), string(want))
}
