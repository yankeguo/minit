package mrunners

import (
	"bytes"
	"context"
	"fmt"
	"github.com/guoyk93/minit/pkg/munit"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func init() {
	Register(munit.KindRender, func(opts RunnerOptions) (runner Runner, err error) {
		if err = opts.Unit.RequireFiles(); err != nil {
			return
		}

		runner.Order = 20
		runner.Func = &runnerRender{RunnerOptions: opts}
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
	tmpl := template.New("__main__").Funcs(Funcs).Option("missingkey=zero")
	if tmpl, err = tmpl.Parse(string(buf)); err != nil {
		err = fmt.Errorf("failed loading %s: %s", name, err.Error())
		return
	}
	out := &bytes.Buffer{}
	if err = tmpl.Execute(out, map[string]interface{}{
		"Env": env,
	}); err != nil {
		err = fmt.Errorf("failed rendering %s: %s", name, err.Error())
		return
	}
	content := out.Bytes()
	if !r.Unit.Raw {
		content = sanitizeLines(content)
	}
	if err = os.WriteFile(name, content, 0755); err != nil {
		err = fmt.Errorf("failed writing %s: %s", name, err.Error())
		return
	}
	return
}

func (r *runnerRender) Do(ctx context.Context) {
	r.Print("controller started")
	defer r.Print("controller exited")

	env := osEnviron()

	for _, filePattern := range r.Unit.Files {
		var err error
		var names []string
		if names, err = filepath.Glob(filePattern); err != nil {
			r.Error(fmt.Sprintf("failed globbing: %s: %s", filePattern, err.Error()))
			continue
		}
		for _, name := range names {
			if err = r.doFile(ctx, name, env); err == nil {
				r.Print("done rendering: " + name)
			} else {
				r.Error("failed rendering: " + name + ": " + err.Error())
			}
		}
	}
}

func osEnviron() map[string]string {
	out := make(map[string]string)
	envs := os.Environ()
	for _, entry := range envs {
		splits := strings.SplitN(entry, "=", 2)
		if len(splits) == 2 {
			out[strings.TrimSpace(splits[0])] = strings.TrimSpace(splits[1])
		}
	}
	return out
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
