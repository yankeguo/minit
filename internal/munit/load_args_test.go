package munit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadArgs(t *testing.T) {
	unit, ok, err := LoadArgs([]string{
		"hello",
		"world",
	})
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, []string{"hello", "world"}, unit.Command)

	unit, ok, err = LoadArgs([]string{
		"minit",
		"/usr/bin/minit",
		"hello",
		"world",
	})
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, []string{"hello", "world"}, unit.Command)

	unit, ok, err = LoadArgs([]string{
		"minit",
		"--a",
		"--b",
		"--",
		"hello",
		"world",
	})
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, []string{"hello", "world"}, unit.Command)

	unit, ok, err = LoadArgs([]string{
		"minit",
		"--a",
		"--b",
		"--",
	})
	require.NoError(t, err)
	require.False(t, ok)

	unit, ok, err = LoadArgs([]string{
		"--a",
		"--b",
		"--",
	})
	require.NoError(t, err)
	require.False(t, ok)

	unit, ok, err = LoadArgs([]string{
		"minit",
		"--once",
		"--b",
		"--",
		"sleep",
		"30",
	})
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, []string{"sleep", "30"}, unit.Command)
	require.Equal(t, KindOnce, unit.Kind)

	unit, ok, err = LoadArgs([]string{
		"--once",
		"--b",
		"--",
		"sleep",
		"30",
	})
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, []string{"sleep", "30"}, unit.Command)
	require.Equal(t, KindOnce, unit.Kind)
}
