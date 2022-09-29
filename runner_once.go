package main

import (
	"context"
	"fmt"
	"github.com/guoyk93/grace/gracelog"
)

const KindOnce = "once"

type OnceRunner struct {
	Unit
	logger *gracelog.ProcLogger
}

func (r *OnceRunner) Run(ctx context.Context) {
	r.logger.Printf("控制器启动")
	defer r.logger.Printf("控制器退出")
	if err := execute(r.ExecuteOptions, r.logger); err != nil {
		r.logger.Errorf("启动失败: %s", err.Error())
		return
	}
}

func NewOnceRunner(unit Unit, logger *gracelog.ProcLogger) (Runner, error) {
	if len(unit.Command) == 0 {
		return nil, fmt.Errorf("没有指定命令，检查 command 字段")
	}
	return &OnceRunner{
		Unit:   unit,
		logger: logger,
	}, nil
}
