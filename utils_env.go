package main

import (
	"os"
	"strconv"
	"strings"
)

func StringEnv(out *string, key, inv string) {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		val = inv
	}
	if *out == "" {
		*out = val
	}
}

func BoolEnv(out *bool, key string) {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	*out, _ = strconv.ParseBool(val)
}
