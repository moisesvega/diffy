package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFlagSet(t *testing.T) {
	users := []string{"first", "second"}
	usersString := strings.Join(users, ",")
	give := []string{
		"--settings",
		// Users
		"--phab_users=" + usersString,
	}

	want := &opts{
		settings:  true,
		phabUsers: users,
	}

	got := &opts{}
	pfs := pflag.NewFlagSet("new", pflag.ExitOnError)
	require.NotPanics(t, func() {
		setFlags(pfs, got)
	})
	err := pfs.Parse(give)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, want, got)
}
