package mrunners

import (
	"context"
	"errors"
	"github.com/guoyk93/gg"
	"github.com/guoyk93/minit/pkg/mexec"
	"github.com/guoyk93/minit/pkg/mlog"
	"github.com/guoyk93/minit/pkg/munit"
	"sync"
)

type Runner struct {
	Order int
	Long  bool
	Func  gg.D10[context.Context]
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

type RunnerFactory = gg.F12[RunnerOptions, Runner, error]

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
