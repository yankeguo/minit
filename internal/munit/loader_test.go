package munit

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLoader(t *testing.T) {
	replaceTestEnv(map[string]string{
		"MINIT_ENABLE":  "@default",
		"MINIT_DISABLE": "task-3,task-5",
		"DEBUG_EVERY":   "10s",
	})
	defer restoreTestEnv()

	ld := NewLoader()
	units, skipped, err := ld.Load(LoadOptions{
		Dir: "testdata",
	})

	require.NoError(t, err)
	require.Len(t, units, 1)
	require.Len(t, skipped, 4)
	require.Equal(t, "task-4", units[0].Name)
	require.Equal(t, "@every 10s", units[0].Cron)
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
	replaceTestEnv(map[string]string{
		"MINIT_MAIN":                 "legacy main",
		"MINIT_UNIT_CACHE_COMMAND":   "redis-server",
		"MINIT_UNIT_INITIAL_COMMAND": "touch /tmp/initial",
		"MINIT_UNIT_INITIAL_NAME":    "job-initial",
		"MINIT_UNIT_INITIAL_ENV":     "ZAA=ZBB;ZCC=ZDD",
		"MINIT_UNIT_INITIAL_KIND":    "once",
	})
	defer restoreTestEnv()

	l := NewLoader()
	units, _, err := l.Load(LoadOptions{Env: true})
	require.NoError(t, err)
	require.Len(t, units, 3)

	sort.Slice(units, func(i, j int) bool {
		return units[i].Name < units[j].Name
	})

	require.Equal(t, "env-cache", units[0].Name)
	require.Equal(t, []string{"redis-server"}, units[0].Command)
	require.Equal(t, KindDaemon, units[0].Kind)

	require.Equal(t, "env-main", units[1].Name)
	require.Equal(t, []string{"legacy", "main"}, units[1].Command)
	require.Equal(t, KindDaemon, units[1].Kind)

	require.Equal(t, "job-initial", units[2].Name)
	require.Equal(t, []string{"touch", "/tmp/initial"}, units[2].Command)
	require.Equal(t, KindOnce, units[2].Kind)
}
