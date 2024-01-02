package menv

import (
	"os"
	"strings"

	"github.com/yankeguo/minit/pkg/mtmpl"
)

const (
	PrefixMinitEnv = "MINIT_ENV_"
)

// Construct create the env map with current system environ, extra and rendering MINIT_ENV_ prefixed keys
func Construct(extra map[string]string) (envs map[string]string, err error) {
	envs = make(map[string]string)
	// system env
	for _, item := range os.Environ() {
		splits := strings.SplitN(item, "=", 2)
		var k, v string
		if len(splits) > 0 {
			k = splits[0]
			if len(splits) > 1 {
				v = splits[1]
			}
			envs[k] = v
		}
	}
	// merge extra env
	Merge(envs, extra)
	// render MINIT_ENV_XXX
	for k, v := range envs {
		if !strings.HasPrefix(k, PrefixMinitEnv) {
			continue
		}
		k = strings.TrimPrefix(k, PrefixMinitEnv)
		var buf []byte
		if buf, err = mtmpl.Execute(v, map[string]any{"Env": envs}); err != nil {
			return
		}
		envs[k] = string(buf)
	}

	return
}
