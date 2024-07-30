package menv

import "os"

var (
	osEnviron = os.Environ
	osGetenv  = os.Getenv
)
