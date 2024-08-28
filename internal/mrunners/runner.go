package mrunners

import (
	"context"
	"errors"
	"sync"
)

// RunnerAction is the interface of runner action
type RunnerAction interface {
	Do(ctx context.Context) (err error)
}

// Runner is the struct of runner
type Runner struct {
	Long   bool
	Action RunnerAction
}

var (
	factories                 = map[string]RunnerFactory{}
	factoriesLock sync.Locker = &sync.Mutex{}
)

// RunnerFactory is the type of runner factory
type RunnerFactory = func(opts RunnerOptions) (Runner, error)

// Register registers a runner factory
func Register(name string, factory RunnerFactory) {
	factoriesLock.Lock()
	defer factoriesLock.Unlock()

	factories[name] = factory
}

// Create creates a runner from options
func Create(opts RunnerOptions) (Runner, error) {
	factoriesLock.Lock()
	defer factoriesLock.Unlock()

	if fac, ok := factories[opts.Unit.Kind]; ok {
		return fac(opts)
	} else {
		return Runner{}, errors.New("unknown runner kind: " + opts.Unit.Kind)
	}
}
