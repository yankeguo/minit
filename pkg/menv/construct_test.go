package menv

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuild(t *testing.T) {
	envs, err := Construct(map[string]string{
		"HOME-":         "NONE",
		"MINIT_ENV_BUF": "{{stringsToUpper \"bbb\"}}",
	})
	require.NoError(t, err)
	require.Equal(t, "", envs["HOME"])
	require.Equal(t, "BBB", envs["BUF"])
}
