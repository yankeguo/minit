package mrunners

import (
	"context"
	"github.com/guoyk93/minit/pkg/munit"
	"github.com/robfig/cron/v3"
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
		runner.Func = &runnerCron{RunnerOptions: opts}
		return
	})
}

type runnerCron struct {
	RunnerOptions
}

func (r *runnerCron) Do(ctx context.Context) {
	r.Logger.Printf("runner started")
	defer r.Logger.Printf("runner exited")

	cr := cron.New(cron.WithLogger(cron.PrintfLogger(r.Logger)))
	_, err := cr.AddFunc(r.Unit.Cron, func() {
		r.Logger.Printf("cron triggered")
		_ = r.Exec.Execute(r.Unit.ExecuteOptions(r.Logger))
		r.Logger.Printf("cron finished")
	})

	if err != nil {
		panic(err)
	}

	cr.Start()

	<-ctx.Done()
	<-cr.Stop().Done()
}
