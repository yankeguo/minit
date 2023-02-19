package ggos

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// MustEnv get environment variable from key, if both value and *out is empty, panic
func MustEnv(key string, out *string) {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		if *out == "" {
			panic(errors.New("missing environment variable '" + key + "'"))
		}
	} else {
		*out = val
	}
}

// BoolEnv get bool environment variable from key
func BoolEnv(key string, out *bool) {
	val := strings.TrimSpace(os.Getenv(key))
	if val != "" {
		*out, _ = strconv.ParseBool(val)
	}
}
