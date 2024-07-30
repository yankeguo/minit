package munit

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	m := map[string]string{
		"MINIT_ENABLE":  "@default",
		"MINIT_DISABLE": "task-3,task-5",
		"DEBUG_EVERY":   "10s",
	}

	units, skipped, err := Load(LoadOptions{
		Dir: "testdata",
		Env: m,
	})

	require.NoError(t, err)
	require.Len(t, units, 1)
	require.Len(t, skipped, 4)
	require.Equal(t, "task-4", units[0].Name)
	require.Equal(t, "@every 10s", units[0].Cron)

	m = map[string]string{
		"MINIT_MAIN":                 "legacy main",
		"MINIT_UNIT_CACHE_COMMAND":   "redis-server",
		"MINIT_UNIT_INITIAL_COMMAND": "touch /tmp/initial",
		"MINIT_UNIT_INITIAL_NAME":    "job-initial",
		"MINIT_UNIT_INITIAL_COUNT":   "3",
		"MINIT_UNIT_INITIAL_ENV":     "ZAA=ZBB;ZCC=ZDD",
		"MINIT_UNIT_INITIAL_KIND":    "once",
		"MINIT_UNIT_WHAT_COMMAND":    "sleep 60",
		"MINIT_UNIT_WHAT_GROUP":      "what",
		"MINIT_DISABLE":              "@what",
	}

	units, skipped, err = Load(LoadOptions{
		Args: []string{
			"-once",
			"--",
			"sleep",
			"60",
		},
		Env: m,
	})
	require.NoError(t, err)
	require.Len(t, units, 6)

	sort.Slice(units, func(i, j int) bool {
		return units[i].Name < units[j].Name
	})

	require.Equal(t, []Unit{
		{
			Kind:    "daemon",
			Name:    "env-what",
			Group:   "what",
			Count:   0,
			Command: []string{"sleep", "60"},
		},
	}, skipped)

	require.Equal(t, []Unit{
		{
			Kind:     "once",
			Name:     "arg-main",
			Group:    "default",
			Count:    1,
			Critical: false,
			Env:      map[string]string{"MINIT_UNIT_NAME": "arg-main", "MINIT_UNIT_SUB_ID": "1"},
			Command:  []string{"sleep", "60"},
		},
		{
			Kind:    "daemon",
			Name:    "env-cache",
			Group:   "default",
			Count:   1,
			Env:     map[string]string{"MINIT_UNIT_NAME": "env-cache", "MINIT_UNIT_SUB_ID": "1"},
			Command: []string{"redis-server"},
		},
		{
			Kind:    "daemon",
			Name:    "env-main",
			Group:   "default",
			Count:   1,
			Env:     map[string]string{"MINIT_UNIT_NAME": "env-main", "MINIT_UNIT_SUB_ID": "1"},
			Command: []string{"legacy", "main"},
		},
		{
			Kind:    "once",
			Name:    "job-initial-1",
			Group:   "default",
			Count:   1,
			Env:     map[string]string{"MINIT_UNIT_NAME": "job-initial-1", "MINIT_UNIT_SUB_ID": "1", "ZAA": "ZBB", "ZCC": "ZDD"},
			Command: []string{"touch", "/tmp/initial"},
		},
		{
			Kind:    "once",
			Name:    "job-initial-2",
			Group:   "default",
			Count:   1,
			Env:     map[string]string{"MINIT_UNIT_NAME": "job-initial-2", "MINIT_UNIT_SUB_ID": "2", "ZAA": "ZBB", "ZCC": "ZDD"},
			Command: []string{"touch", "/tmp/initial"},
		},
		{
			Kind:     "once",
			Name:     "job-initial-3",
			Group:    "default",
			Count:    1,
			Critical: false,
			Dir:      "",
			Shell:    "",
			Env:      map[string]string{"MINIT_UNIT_NAME": "job-initial-3", "MINIT_UNIT_SUB_ID": "3", "ZAA": "ZBB", "ZCC": "ZDD"},
			Command:  []string{"touch", "/tmp/initial"},
		},
	}, units)
}

func TestDupOrMakeMap(t *testing.T) {
	var o map[string]any
	duplicateMap(&o)
	require.NotNil(t, o)

	m1a := map[string]string{
		"a": "b",
	}
	m1b := m1a
	duplicateMap(&m1a)
	m1a["c"] = "d"
	require.Equal(t, "d", m1a["c"])
	require.Equal(t, "", m1b["c"])
}
