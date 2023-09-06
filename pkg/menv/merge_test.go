package menv

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMerge(t *testing.T) {
	m := map[string]string{
		"a": "b",
		"c": "d",
	}
	m2 := map[string]string{
		"a-": "",
		"c":  "e",
		"h":  "j",
	}
	Merge(m, m2)
	require.Equal(t, 2, len(m))
	require.Equal(t, "e", m["c"])
	require.Equal(t, "j", m["h"])
}
