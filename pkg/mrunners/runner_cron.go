package mrunners

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/yankeguo/minit/pkg/munit"
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

func (r *runnerCron) Do(ctx context.Context) {
	r.Print("controller started")
	defer r.Print("controller exited")

	if r.Unit.Immediate {
		if err := r.Exec.Execute(r.Unit.ExecuteOptions(r.Logger)); err != nil {
			r.Error("failed executing: " + err.Error())
		}
	}

	cr := cron.New(cron.WithLogger(cron.PrintfLogger(r.Logger)))
	_, err := cr.AddFunc(r.Unit.Cron, func() {
		r.Print("triggered")
		if err := r.Exec.Execute(r.Unit.ExecuteOptions(r.Logger)); err != nil {
			r.Error("failed executing: " + err.Error())
		}
	})

	if err != nil {
		panic(err)
	}

	cr.Start()

	<-ctx.Done()
	<-cr.Stop().Done()
}
