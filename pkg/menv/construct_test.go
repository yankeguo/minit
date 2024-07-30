package menv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func replaceTestEnv(m map[string]string) {
	osGetenv = func(key string) string {
		return m[key]
	}
	osEnviron = func() []string {
		var env []string
		for k, v := range m {
			env = append(env, k+"="+v)
		}
		return env
	}
}

func restoreTestEnv() {
	osGetenv = os.Getenv
	osEnviron = os.Environ
}

func TestBuild(t *testing.T) {
	replaceTestEnv(map[string]string{
		"HOME": "/home/minit",
	})
	defer restoreTestEnv()

	envs, err := Construct(map[string]string{
		"HOME-":         "NONE",
		"MINIT_ENV_BUF": "{{stringsToUpper \"bbb\"}}",
	})
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"BUF": "BBB",
	}, envs)
}
