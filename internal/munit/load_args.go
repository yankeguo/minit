package munit

import (
	"path/filepath"
	"strings"
)

// LoadArgs loads unit from command line arguments
func LoadArgs(args []string) (unit Unit, ok bool, err error) {
	var opts []string

	// fix a history issue
	for len(args) > 0 {
		if filepath.Base(args[0]) == "minit" {
			args = args[1:]
		} else {
			break
		}
	}

	// extract arguments after '--' if existed
	for i, item := range args {
		if item == "--" {
			opts = args[0:i]
			args = args[i+1:]
			break
		}
	}

	if len(args) == 0 {
		return
	}

	unit = Unit{
		Name:    "arg-main",
		Kind:    KindDaemon,
		Command: args,
	}

	// opts decoding
	for _, opt := range opts {
		if strings.HasSuffix(opt, "-"+KindOnce) {
			unit.Kind = KindOnce
		}
	}

	ok = true

	return
}
