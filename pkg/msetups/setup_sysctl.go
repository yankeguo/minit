//go:build linux

package msetups

import (
	"fmt"
	"github.com/guoyk93/minit/pkg/mlog"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	Register(20, setupSysctl)
}

func setupSysctl(logger mlog.ProcLogger) (err error) {
	items := strings.Split(os.Getenv("MINIT_SYSCTL"), ",")
	for _, item := range items {
		splits := strings.SplitN(item, "=", 2)
		if len(splits) != 2 {
			continue
		}

		k, v := strings.TrimSpace(splits[0]), strings.TrimSpace(splits[1])
		if k == "" {
			continue
		}

		filename := filepath.Join(
			append(
				[]string{"/proc", "sys"},
				strings.Split(k, ".")...,
			)...,
		)

		logger.Printf("writing sysctl %s=%s", k, v)

		if err = os.WriteFile(filename, []byte(v), 0644); err != nil {
			err = fmt.Errorf("failed writing sysctl %s=%s: %s", k, v, err.Error())
			return
		}
	}
	return
}
