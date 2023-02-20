package munit

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
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
