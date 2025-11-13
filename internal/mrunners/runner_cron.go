package mrunners

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/yankeguo/minit/internal/munit"
	"github.com/yankeguo/rg"
)

func init() {
	Register(munit.KindCron, func(opts RunnerOptions) (runner Runner, err error) {
		defer rg.Guard(&err)
		rg.Must0(opts.Unit.RequireCommand())
		rg.Must0(opts.Unit.RequireCron())

		// Validate cron expression with detailed error context
		if _, parseErr := cron.ParseStandard(opts.Unit.Cron); parseErr != nil {
			err = fmt.Errorf("cron unit '%s': invalid cron expression '%s': %w", opts.Unit.Name, opts.Unit.Cron, parseErr)
			return
		}

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
