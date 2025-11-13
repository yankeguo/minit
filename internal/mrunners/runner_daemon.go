package mrunners

import (
	"context"
	"time"

	"github.com/yankeguo/minit/internal/munit"
	"github.com/yankeguo/rg"
)

func init() {
	Register(munit.KindDaemon, func(opts RunnerOptions) (runner Runner, err error) {
		defer rg.Guard(&err)
		rg.Must0(opts.Unit.RequireCommand())

		runner.Long = true
		runner.Action = &actionDaemon{RunnerOptions: opts}
		return
	})
}

type actionDaemon struct {
	RunnerOptions
}

func (r *actionDaemon) Do(ctx context.Context) (err error) {
	r.Print("controller started")
	defer r.Print("controller exited")
	defer rg.Guard(&err)

forLoop:
	for {
		if ctx.Err() != nil {
			break forLoop
		}

		err = r.PanicOnCritical("failed executing", r.Execute())

		if ctx.Err() != nil {
			break forLoop
		}

		r.Print("restarting")

		// Create timer for restart delay with proper cleanup
		timer := time.NewTimer(time.Second * 5)
		select {
		case <-timer.C:
			// Timer expired naturally
		case <-ctx.Done():
			// Context cancelled, stop timer to prevent resource leak
			if !timer.Stop() {
				// Timer already fired, drain the channel
				<-timer.C
			}
			break forLoop
		}
	}

	return
}
