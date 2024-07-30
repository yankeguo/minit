package munit

import "os"

var (
	osEnviron   = os.Environ
	osGetenv    = os.Getenv
	osExpandEnv = os.ExpandEnv
)
