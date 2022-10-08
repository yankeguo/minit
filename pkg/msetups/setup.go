package msetups

import (
	"github.com/guoyk93/gg"
	"github.com/guoyk93/minit/pkg/mlog"
	"sort"
	"sync"
)

type SetupFunc = gg.F11[mlog.ProcLogger, error]

var (
	setupsLock sync.Locker = &sync.Mutex{}
	setups     []gg.T2[int, SetupFunc]
)

func Register(order int, fn SetupFunc) {
	setupsLock.Lock()
	defer setupsLock.Unlock()

	setups = append(setups, gg.T2[int, SetupFunc]{A: order, B: fn})
}

func Setup(logger mlog.ProcLogger) (err error) {
	setupsLock.Lock()
	defer setupsLock.Unlock()

	sort.Slice(setups, func(i, j int) bool {
		return setups[i].A > setups[j].A
	})

	for _, setup := range setups {
		if err = setup.B(logger); err != nil {
			return
		}
	}

	return
}
