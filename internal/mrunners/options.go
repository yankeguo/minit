package mrunners

import (
	"github.com/yankeguo/minit/internal/mexec"
	"github.com/yankeguo/minit/internal/mlog"
	"github.com/yankeguo/minit/internal/munit"
)

type RunnerOptions struct {
	Unit   munit.Unit
	Exec   mexec.Manager
	Logger mlog.ProcLogger
}

func (ro RunnerOptions) Print(message string) {
	ro.Logger.Print("minit: " + ro.Unit.Kind + "/" + ro.Unit.Name + ": " + message)
}

func (ro RunnerOptions) Printf(layout string, items ...any) {
	ro.Logger.Printf("minit: "+ro.Unit.Kind+"/"+ro.Unit.Name+": "+layout, items...)
}

func (ro RunnerOptions) Error(message string) {
	ro.Logger.Error("minit: " + ro.Unit.Kind + "/" + ro.Unit.Name + ": " + message)
}

func (ro RunnerOptions) Errorf(layout string, items ...any) {
	ro.Logger.Errorf("minit: "+ro.Unit.Kind+"/"+ro.Unit.Name+": "+layout, items...)
}

func (ro RunnerOptions) Execute() error {
	return ro.Exec.Execute(ro.Unit.ExecuteOptions(ro.Logger))
}
