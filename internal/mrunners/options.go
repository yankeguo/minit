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

func (ro RunnerOptions) PanicOnCritical(message string, err error) error {
	if err == nil {
		return nil
	}

	ro.Error(message + ": " + err.Error())

	if ro.Unit.Critical {
		panic(err)
	} else {
		return nil
	}
}
