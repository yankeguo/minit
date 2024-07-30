package menv

import (
	"strings"

	"github.com/yankeguo/minit/pkg/mtmpl"
)

const (
	EnvPrefixEnv = "MINIT_ENV_"
)

// Construct create the env map with current system environ, extra and rendering MINIT_ENV_ prefixed keys
func Construct(extra map[string]string) (envs map[string]string, err error) {
	envs = make(map[string]string)

	// system env
	for _, item := range osEnviron() {
		splits := strings.SplitN(item, "=", 2)
		var key, val string
		if len(splits) > 0 {
			key = splits[0]
			if len(splits) > 1 {
				val = splits[1]
			}
			envs[key] = val
		}
	}

	// merge extra env
	Merge(envs, extra)

	// render MINIT_ENV_XXX
	for key, val := range envs {
		if !strings.HasPrefix(key, EnvPrefixEnv) {
			continue
		}
		effectiveKey := strings.TrimPrefix(key, EnvPrefixEnv)
		var buf []byte
		if buf, err = mtmpl.Execute(val, map[string]any{"Env": envs}); err != nil {
			return
		}
		// set the rendered value
		envs[effectiveKey] = string(buf)
		// remove the original key
		delete(envs, key)
	}

	return
}
