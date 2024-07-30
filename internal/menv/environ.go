package menv

import (
	"os"
	"strings"
)

// Environ returns the system environment variables as a map
func Environ() (m map[string]string) {
	m = make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.SplitN(item, "=", 2)
		var key, val string
		if len(splits) > 0 {
			key = splits[0]
			if len(splits) > 1 {
				val = splits[1]
			}
			m[key] = val
		}
	}
	return
}
