package mrunners

import (
	"context"
	"errors"
	"sync"

	"github.com/yankeguo/minit/pkg/mexec"
	"github.com/yankeguo/minit/pkg/mlog"
	"github.com/yankeguo/minit/pkg/munit"
)

type RunnerAction interface {
	Do(ctx context.Context)
}

type Runner struct {
	Order  int
	Long   bool
	Action RunnerAction
}

var (
	factories                 = map[string]RunnerFactory{}
	factoriesLock sync.Locker = &sync.Mutex{}
)

type RunnerOptions struct {
	Unit   munit.Unit
	Exec   mexec.Manager
	Logger mlog.ProcLogger
}

func (ro RunnerOptions) Print(message string) {
	ro.Logger.Print("minit: " + ro.Unit.Kind + "/" + ro.Unit.Name + ": " + message)
}

func (ro RunnerOptions) Error(message string) {
	ro.Logger.Error("minit: " + ro.Unit.Kind + "/" + ro.Unit.Name + ": " + message)
}

type RunnerFactory = func(opts RunnerOptions) (Runner, error)

func Register(name string, factory RunnerFactory) {
	factoriesLock.Lock()
	defer factoriesLock.Unlock()

	factories[name] = factory
}

func Create(opts RunnerOptions) (Runner, error) {
	factoriesLock.Lock()
	defer factoriesLock.Unlock()

	if fac, ok := factories[opts.Unit.Kind]; ok {
		return fac(opts)
	} else {
		return Runner{}, errors.New("unknown runner kind: " + opts.Unit.Kind)
	}
}
