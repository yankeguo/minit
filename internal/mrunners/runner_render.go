package mrunners

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"unicode"

	"github.com/yankeguo/minit/internal/menv"
	"github.com/yankeguo/minit/internal/mtmpl"
	"github.com/yankeguo/minit/internal/munit"
)

func init() {
	Register(munit.KindRender, func(opts RunnerOptions) (runner Runner, err error) {
		if err = opts.Unit.RequireFiles(); err != nil {
			return
		}

		runner.Order = 10
		runner.Action = &actionRender{RunnerOptions: opts}
		return
	})
}

type actionRender struct {
	RunnerOptions
}

func (r *actionRender) doFile(ctx context.Context, name string, env map[string]string) (err error) {
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

func (r *actionRender) Do(ctx context.Context) (err error) {
	r.Print("controller started")
	defer r.Print("controller exited")

	var env map[string]string

	if env, err = menv.Construct(r.Unit.Env); err != nil {
		r.Error("failed constructing environments variables: " + err.Error())
		return
	}

	allNames := map[string]struct{}{}

	for _, filePattern := range r.Unit.Files {
		var names []string

		if names, err = filepath.Glob(filePattern); err != nil {
			r.Errorf("failed globbing: %s: %s", filePattern, err.Error())

			if r.Unit.Critical {
				return
			} else {
				err = nil
			}

			continue
		}

		for _, name := range names {
			allNames[name] = struct{}{}
		}
	}

	for name := range allNames {
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

	return
}

// sanitizeLines removes empty lines and trailing spaces
func sanitizeLines(s []byte) []byte {
	var out [][]byte
	for _, line := range bytes.Split(s, []byte("\n")) {
		line = bytes.TrimRightFunc(line, unicode.IsSpace)
		if len(line) == 0 {
			continue
		}
		out = append(out, line)
	}
	return bytes.Join(out, []byte("\n"))
}
