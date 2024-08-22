//go:build linux

package msetups

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/yankeguo/minit/internal/mlog"
)

const (
	controlFileTHP = "/sys/kernel/mm/transparent_hugepage/enabled"
)

func init() {
	Register(40, setupTHP)
}

func setupTHP(logger mlog.ProcLogger) (err error) {
	val := strings.TrimSpace(os.Getenv("MINIT_THP"))
	if val == "" {
		return
	}
	var buf []byte
	if buf, err = os.ReadFile(controlFileTHP); err != nil {
		err = fmt.Errorf("failed reading THP configuration %s: %s", controlFileTHP, err.Error())
		return
	}
	logger.Printf("current THP configuration: %s", bytes.TrimSpace(buf))
	logger.Printf("writing THP configuration: %s", val)
	if err = os.WriteFile(controlFileTHP, []byte(val), 0644); err != nil {
		err = fmt.Errorf("failed writing THP configuration %s: %s", controlFileTHP, err.Error())
		return
	}
	if buf, err = os.ReadFile(controlFileTHP); err != nil {
		err = fmt.Errorf("failed reading THP configuration %s: %s", controlFileTHP, err.Error())
		return
	}
	logger.Printf("current THP configuration: %s", bytes.TrimSpace(buf))
	return
}
