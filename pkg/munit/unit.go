package munit

import (
	"github.com/guoyk93/minit/pkg/mexec"
	"github.com/guoyk93/minit/pkg/mlog"
)

const (
	DefaultGroup = "default"
)

const (
	KindDaemon = "daemon"
	KindOnce   = "once"
	KindCron   = "cron"
	KindRender = "render"
)

type Unit struct {
	Kind  string `yaml:"kind"`  // kind of unit
	Name  string `yaml:"name"`  // name of unit
	Group string `yaml:"group"` // group of unit
	Count int    `yaml:"count"` // replicas of unit

	// execution options
	Dir     string            `yaml:"dir"`
	Shell   string            `yaml:"shell"`
	Env     map[string]string `yaml:"env"`
	Command []string          `yaml:"command"`
	Charset string            `yaml:"charset"`

	// for 'render' only
	Raw   bool     `yaml:"raw"`   // don't trim white spaces for 'render'
	Files []string `yaml:"files"` // files to process

	// for 'cron' only
	Cron string `yaml:"cron"` // cron syntax
}

func (u Unit) ExecuteOptions(logger mlog.ProcLogger) mexec.ExecuteOptions {
	return mexec.ExecuteOptions{
		Dir:     u.Dir,
		Shell:   u.Shell,
		Env:     u.Env,
		Command: u.Command,
		Charset: u.Charset,

		Logger:          logger,
		IgnoreExecError: true,
	}
}

func (u Unit) CanonicalName() string {
	return u.Kind + "/" + u.Name
}
