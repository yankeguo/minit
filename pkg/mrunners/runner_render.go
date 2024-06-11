package mrunners

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yankeguo/minit/pkg/menv"
	"github.com/yankeguo/minit/pkg/mtmpl"
	"github.com/yankeguo/minit/pkg/munit"
)

func init() {
	Register(munit.KindRender, func(opts RunnerOptions) (runner Runner, err error) {
		if err = opts.Unit.RequireFiles(); err != nil {
			return
		}

		runner.Order = 10
		runner.Action = &runnerRender{RunnerOptions: opts}
		return
	})
}

type runnerRender struct {
	RunnerOptions
}

func (r *runnerRender) doFile(ctx context.Context, name string, env map[string]string) (err error) {
	var buf []byte
	if buf, err = os.ReadFile(name); err != nil {
		err = fmt.Errorf("failed reading %s: %s", name, err.Error())
		return
	}
	var content []byte
	if content, err = mtmpl.Execute(string(buf), map[string]any{
		"Env": env,
	}); err != nil {
		err = fmt.Errorf("failed rendering %s: %s", name, err.Error())
		return
	}
	if !r.Unit.Raw {
		content = sanitizeLines(content)
	}
	if err = os.WriteFile(name, content, 0755); err != nil {
		err = fmt.Errorf("failed writing %s: %s", name, err.Error())
		return
	}
	return
}

func (r *runnerRender) Do(ctx context.Context) (err error) {
	r.Print("controller started")
	defer r.Print("controller exited")

	var env map[string]string

	if env, err = menv.Construct(r.Unit.Env); err != nil {
		r.Error("failed constructing environments variables: " + err.Error())
		return
	}

	for _, filePattern := range r.Unit.Files {
		var names []string

		if names, err = filepath.Glob(filePattern); err != nil {
			r.Error(fmt.Sprintf("failed globbing: %s: %s", filePattern, err.Error()))

			if r.Unit.Critical {
				return
			} else {
				err = nil
			}

			continue
		}

		for _, name := range names {
			if err = r.doFile(ctx, name, env); err != nil {
				r.Error("failed rendering: " + name + ": " + err.Error())

				if r.Unit.Critical {
					return
				} else {
					err = nil
				}

				continue
			}

			r.Print("done rendering: " + name)
		}
	}

	return
}

func sanitizeLines(s []byte) []byte {
	lines := bytes.Split(s, []byte("\n"))
	out := &bytes.Buffer{}
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		out.Write(line)
		out.WriteRune('\n')
	}
	return out.Bytes()
}
