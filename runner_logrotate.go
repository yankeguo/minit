package main

import (
	"context"
	"github.com/guoyk93/grace/gracelog"
)

const KindLogrotate = "logrotate"

type LogrotateRunner struct {
	Unit
	logger *gracelog.ProcLogger
}

func (l *LogrotateRunner) Run(ctx context.Context) {
	l.logger.Printf("控制器启动")
	defer l.logger.Printf("控制器退出")

	l.logger.Error("警告：minit 的 logrotate 功能从未完成开发，当前已经被弃用")

	<-ctx.Done()
}

func NewLogrotateRunner(unit Unit, logger *gracelog.ProcLogger) (Runner, error) {
	return &LogrotateRunner{
		Unit:   unit,
		logger: logger,
	}, nil
}
