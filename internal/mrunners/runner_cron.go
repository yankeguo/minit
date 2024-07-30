package mrunners

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/yankeguo/minit/internal/munit"
	"github.com/yankeguo/rg"
)

func init() {
	Register(munit.KindCron, func(opts RunnerOptions) (runner Runner, err error) {
		if err = opts.Unit.RequireCommand(); err != nil {
			return
		}
		if err = opts.Unit.RequireCron(); err != nil {
			return
		}
		if _, err = cron.ParseStandard(opts.Unit.Cron); err != nil {
			return
		}

		runner.Order = 30
		runner.Long = true
		runner.Action = &actionCron{RunnerOptions: opts}
		return
	})
}

type actionCron struct {
	RunnerOptions
}

func (r *actionCron) Do(ctx context.Context) (err error) {
	r.Print("controller started")
	defer r.Print("controller exited")
	defer rg.Guard(&err)

	if r.Unit.Immediate {
		err = r.PanicOnCritical("failed executing", r.Execute())
	}

	cr := cron.New(cron.WithLogger(cron.PrintfLogger(r.Logger)))

	var chErr chan error

	if r.Unit.Critical {
		chErr = make(chan error, 1)
	}

	rg.Must(
		cr.AddFunc(
			r.Unit.Cron,
			func() {
				r.Print("triggered")
				if err := func() (err error) {
					defer rg.Guard(&err)
					return r.PanicOnCritical("failed executing", r.Execute())
				}(); err != nil {
					if chErr != nil {
						select {
						case chErr <- err:
						default:
						}
					}
				}
			},
		),
	)

	cr.Start()

	if chErr != nil {
		select {
		case <-ctx.Done():
		case err = <-chErr:
		}
	} else {
		<-ctx.Done()
	}

	<-cr.Stop().Done()

	return
}
