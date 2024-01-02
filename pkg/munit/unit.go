package munit

import (
	"errors"

	"github.com/yankeguo/minit/pkg/mexec"
	"github.com/yankeguo/minit/pkg/mlog"
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

var (
	knownUnitKind = map[string]struct{}{
		KindDaemon: {},
		KindOnce:   {},
		KindCron:   {},
		KindRender: {},
	}
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
	Cron      string `yaml:"cron"` // cron syntax
	Immediate bool   `yaml:"immediate"`
}

func (u Unit) RequireCommand() error {
	if len(u.Command) == 0 {
		return errors.New("missing unit field: command")
	}
	return nil
}

func (u Unit) RequireFiles() error {
	if len(u.Files) == 0 {
		return errors.New("missing unit field: command")
	}
	return nil
}

func (u Unit) RequireCron() error {
	if len(u.Cron) == 0 {
		return errors.New("missing unit field: cron")
	}
	return nil
}

func (u Unit) ExecuteOptions(logger mlog.ProcLogger) mexec.ExecuteOptions {
	return mexec.ExecuteOptions{
		Name: u.Kind + "/" + u.Name,

		Dir:     u.Dir,
		Shell:   u.Shell,
		Env:     u.Env,
		Command: u.Command,
		Charset: u.Charset,

		Logger:          logger,
		IgnoreExecError: true,
	}
}
