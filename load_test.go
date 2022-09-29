package main

import (
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestLoadDir(t *testing.T) {
	units, err := LoadDir(filepath.Join("testdata", "minit.d"))
	require.NoError(t, err)
	require.Equal(t, "cron", units[7].Kind)
	require.Equal(t, "@every 10s", units[7].Cron)
}
