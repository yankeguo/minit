package main

import (
	"context"
	"fmt"
	"github.com/guoyk93/minit/pkg/mlog"
	"github.com/guoyk93/minit/pkg/munit"
)

const KindOnce = "once"

type OnceRunner struct {
	munit.Unit
	logger mlog.ProcLogger
}

func (r *OnceRunner) Run(ctx context.Context) {
	r.logger.Printf("控制器启动")
	defer r.logger.Printf("控制器退出")
	if err := EXE.Execute(r.ExecuteOptions(r.logger)); err != nil {
		r.logger.Errorf("启动失败: %s", err.Error())
		return
	}
}

func NewOnceRunner(unit munit.Unit, logger mlog.ProcLogger) (Runner, error) {
	if len(unit.Command) == 0 {
		return nil, fmt.Errorf("没有指定命令，检查 command 字段")
	}
	return &OnceRunner{
		Unit:   unit,
		logger: logger,
	}, nil
}
