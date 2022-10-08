//go:build !linux

package main

import (
	"github.com/guoyk93/minit/pkg/mlog"
)

func setupTHP() error {
	return nil
}

func setupSysctl() error {
	return nil
}

func setupRLimits() error {
	return nil
}

func setupZombies(log mlog.ProcLogger) {
	return
}
