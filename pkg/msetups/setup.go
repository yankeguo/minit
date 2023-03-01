package msetups

import (
	"github.com/guoyk93/minit/pkg/mlog"
	"sort"
	"sync"
)

type SetupFunc = func(log mlog.ProcLogger) error

type setupItem struct {
	order int
	fn    SetupFunc
}

var (
	setupsLock sync.Locker = &sync.Mutex{}
	setups     []setupItem
)

func Register(order int, fn SetupFunc) {
	setupsLock.Lock()
	defer setupsLock.Unlock()

	setups = append(setups, setupItem{order: order, fn: fn})
}

func Setup(logger mlog.ProcLogger) (err error) {
	setupsLock.Lock()
	defer setupsLock.Unlock()

	sort.Slice(setups, func(i, j int) bool {
		return setups[i].order > setups[j].order
	})

	for _, setup := range setups {
		if err = setup.fn(logger); err != nil {
			return
		}
	}

	return
}
