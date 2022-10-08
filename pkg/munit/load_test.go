package munit

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
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
}

func TestLoadEnv(t *testing.T) {
	os.Setenv("MINIT_MAIN", "hello 'world destroyer'")
	os.Setenv("MINIT_MAIN_KIND", "cron")
	os.Setenv("MINIT_MAIN_NAME", "test-main")
	os.Setenv("MINIT_MAIN_CRON", "1 2 3 4 5")
	os.Setenv("MINIT_MAIN_GROUP", "bbb")
	os.Setenv("MINIT_MAIN_CHARSET", "gbk")

	unit, ok, err := LoadEnv()
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, Unit{
		Kind:  "cron",
		Name:  "test-main",
		Cron:  "1 2 3 4 5",
		Group: "bbb",
		Command: []string{
			"hello",
			"world destroyer",
		},
		Charset: "gbk",
	}, unit)
}
