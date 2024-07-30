package munit

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDir(t *testing.T) {
	units, err := LoadDir("testdata")
	require.NoError(t, err)
	require.NotEmpty(t, units)

	sort.Slice(units, func(i, j int) bool {
		return units[i].Name < units[j].Name
	})

	require.Equal(t, []Unit{
		{
			Kind:    "once",
			Name:    "task-1",
			Group:   "group-echo",
			Command: []string{"echo", "once", "$HOME"},
		},
		{
			Kind:         "daemon",
			Name:         "task-2",
			Group:        "group-echo",
			Count:        3,
			Critical:     true,
			Shell:        "/bin/bash",
			Command:      []string{"sleep 1 && echo hello world"},
			Charset:      "gbk",
			SuccessCodes: []int{0, 1, 2},
		},
		{
			Kind:    "daemon",
			Name:    "task-3",
			Count:   3,
			Command: []string{"sleep", "5"},
		},
		{
			Kind:      "cron",
			Name:      "task-4",
			Command:   []string{"echo", "cron"},
			Cron:      "@every ${DEBUG_EVERY}",
			Immediate: true,
		},
		{
			Kind:  "render",
			Name:  "task-5",
			Raw:   true,
			Files: []string{"testdata/conf/*.conf"},
		},
	}, units)
}
