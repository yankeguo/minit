package menv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstruct(t *testing.T) {
	envs, err := Construct(map[string]string{
		"HOME": "/home/minit",
	}, map[string]string{
		"HOME-":         "NONE",
		"MINIT_ENV_BUF": "{{stringsToUpper \"bbb\"}}",
	})
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"BUF": "BBB",
	}, envs)
}
