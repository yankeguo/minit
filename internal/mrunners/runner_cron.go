package mrunners

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/yankeguo/minit/internal/munit"
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
		runner.Action = &runnerCron{RunnerOptions: opts}
		return
	})
}

type runnerCron struct {
	RunnerOptions
}

func (r *runnerCron) Do(ctx context.Context) (err error) {
	r.Print("controller started")
	defer r.Print("controller exited")

	if r.Unit.Immediate {
		if err = r.Exec.Execute(r.Unit.ExecuteOptions(r.Logger)); err != nil {
			r.Error("failed executing: " + err.Error())
			if r.Unit.Critical {
				return
			} else {
				err = nil
			}
		}
	}

	cr := cron.New(cron.WithLogger(cron.PrintfLogger(r.Logger)))

	var chErr chan error

	if r.Unit.Critical {
		chErr = make(chan error, 1)
	}

	if _, err = cr.AddFunc(r.Unit.Cron, func() {
		r.Print("triggered")
		if err := r.Exec.Execute(r.Unit.ExecuteOptions(r.Logger)); err != nil {
			r.Error("failed executing: " + err.Error())
			if chErr != nil {
				select {
				case chErr <- err:
				default:
				}
			} else {
				err = nil
			}
		}
	}); err != nil {
		// should fail since we have checked in init
		return
	}

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
