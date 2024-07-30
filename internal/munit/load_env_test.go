package munit

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectEnvInfixes(t *testing.T) {
	env := map[string]string{
		"MINIT_UNIT_MAIN_COMMAND": "hello",
		"MINIT_UNIT_COMMAND":      "hello",
		"MINIT_UNIT_A_FILES":      "hello",
		"MINIT_UNIT_B_FILES":      "hello",
		"MINIT_UNIT_B_KIND":       "render",
	}

	infixes := DetectEnvInfixes(env)
	sort.StringSlice(infixes).Sort()

	require.Equal(t, []string{"B", "MAIN"}, infixes)
}

func TestLoadEnvWithInfix(t *testing.T) {
	env := map[string]string{
		"MINIT_UNIT_A1_COMMAND":       "echo 'hello world'",
		"MINIT_UNIT_A2_COMMAND":       "echo 'hello world'",
		"MINIT_UNIT_A2_KIND":          "cron",
		"MINIT_UNIT_A2_CRON":          "* * * * *",
		"MINIT_UNIT_A2_NAME":          "a2",
		"MINIT_UNIT_A2_IMMEDIATE":     "true",
		"MINIT_UNIT_A2_GROUP":         "abc",
		"MINIT_UNIT_A2_COUNT":         "3",
		"MINIT_UNIT_A2_DIR":           "/opt",
		"MINIT_UNIT_A2_SHELL":         "/bin/zsh",
		"MINIT_UNIT_A2_CHARSET":       "gbk",
		"MINIT_UNIT_A2_ENV":           "a=b;c=d",
		"MINIT_UNIT_A2_CRITICAL":      "true",
		"MINIT_UNIT_A2_SUCCESS_CODES": "114,514",
		"MINIT_UNIT_A3_COMMAND":       "echo 'hello world'",
		"MINIT_UNIT_A3_KIND":          "once",
		"MINIT_UNIT_A3_BLOCKING":      "false",
		"MINIT_UNIT_A4_KIND":          "render",
		"MINIT_UNIT_A4_FILES":         "hello.txt;world.txt",
		"MINIT_UNIT_A4_RAW":           "true",
	}

	unit, ok, err := LoadEnvWithInfix(env, "A1")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, Unit{
		Kind: KindDaemon,
		Name: "env-a1",
		Command: []string{
			"echo",
			"hello world",
		},
	}, unit)

	unit, ok, err = LoadEnvWithInfix(env, "A2")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, Unit{
		Kind: KindCron,
		Name: "a2",
		Command: []string{
			"echo",
			"hello world",
		},
		Immediate: true,
		Cron:      "* * * * *",
		Group:     "abc",
		Count:     3,
		Dir:       "/opt",
		Shell:     "/bin/zsh",
		Charset:   "gbk",
		Env: map[string]string{
			"a": "b",
			"c": "d",
		},
		Critical:     true,
		SuccessCodes: []int{114, 514},
	}, unit)

	blockingTrue := false

	unit, ok, err = LoadEnvWithInfix(env, "A3")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, Unit{
		Kind:     KindOnce,
		Name:     "env-a3",
		Blocking: &blockingTrue,
		Command: []string{
			"echo",
			"hello world",
		},
	}, unit)

	unit, ok, err = LoadEnvWithInfix(env, "A4")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, Unit{
		Kind: KindRender,
		Name: "env-a4",
		Raw:  true,
		Files: []string{
			"hello.txt",
			"world.txt",
		},
	}, unit)
}

func TestLoadEnv(t *testing.T) {
	env := map[string]string{
		"MINIT_MAIN":         "hello 'world destroyer'",
		"MINIT_MAIN_KIND":    "cron",
		"MINIT_MAIN_NAME":    "test-main",
		"MINIT_MAIN_CRON":    "1 2 3 4 5",
		"MINIT_MAIN_GROUP":   "bbb",
		"MINIT_MAIN_CHARSET": "gbk",
	}

	unit, ok, err := LoadEnv(env)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, Unit{
		Kind:  "cron",
		Name:  "test-main",
		Cron:  "1 2 3 4 5",
		Group: "bbb",
		Command: []string{
			"hello",
			"world destroyer",
		},
		Charset: "gbk",
	}, unit)
}
