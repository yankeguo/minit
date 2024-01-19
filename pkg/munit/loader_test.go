package munit

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLoader(t *testing.T) {
	os.Setenv("MINIT_ENABLE", "@default")
	os.Setenv("MINIT_DISABLE", "task-3,task-5")
	ld := NewLoader()
	units, skipped, err := ld.Load(LoadOptions{
		Dir: "testdata",
	})

	require.NoError(t, err)
	require.Len(t, units, 1)
	require.Len(t, skipped, 4)
	require.Equal(t, "task-4", units[0].Name)
}

func TestDupOrMakeMap(t *testing.T) {
	var o map[string]any
	dupOrMakeMap(&o)
	require.NotNil(t, o)

	m1a := map[string]string{
		"a": "b",
	}
	m1b := m1a
	dupOrMakeMap(&m1a)
	m1a["c"] = "d"
	require.Equal(t, "d", m1a["c"])
	require.Equal(t, "", m1b["c"])
}

func TestLoaderWithNewEnv(t *testing.T) {
	os.Unsetenv("MINIT_ENABLE")
	os.Unsetenv("MINIT_DISABLE")
	os.Unsetenv("MINIT_MAIN_NAME")
	os.Unsetenv("MINIT_MAIN_KIND")
	os.Unsetenv("MINIT_MAIN_ONCE")
	os.Setenv("MINIT_MAIN", "legacy main")
	os.Setenv("MINIT_UNIT_CACHE_COMMAND", "redis-server")
	os.Setenv("MINIT_UNIT_INITIAL_COMMAND", "touch /tmp/initial")
	os.Setenv("MINIT_UNIT_INITIAL_NAME", "job-initial")
	os.Setenv("MINIT_UNIT_INITIAL_ENV", "ZAA=ZBB;ZCC=ZDD")
	os.Setenv("MINIT_UNIT_INITIAL_KIND", "once")

	l := NewLoader()
	units, _, err := l.Load(LoadOptions{Env: true})
	require.NoError(t, err)
	require.Len(t, units, 3)

	require.Equal(t, "env-main", units[0].Name)
	require.Equal(t, []string{"legacy", "main"}, units[0].Command)
	require.Equal(t, KindDaemon, units[0].Kind)

	require.Equal(t, "env-cache", units[1].Name)
	require.Equal(t, []string{"redis-server"}, units[1].Command)
	require.Equal(t, KindDaemon, units[1].Kind)

	require.Equal(t, "job-initial", units[2].Name)
	require.Equal(t, []string{"touch", "/tmp/initial"}, units[2].Command)
	require.Equal(t, KindOnce, units[2].Kind)
}
