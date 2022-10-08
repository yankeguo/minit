package mrunners

import (
	"bytes"
	"context"
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

func (r *runnerRender) Do(ctx context.Context) {
	r.Logger.Printf("runner started")
	defer r.Logger.Printf("runner exited")

	env := osEnviron()

	for _, filePattern := range r.Unit.Files {
		var err error
		var names []string
		if names, err = filepath.Glob(filePattern); err != nil {
			r.Logger.Errorf("glob %s failed: %s", filePattern, err.Error())
			continue
		}
		for _, name := range names {
			var buf []byte
			if buf, err = os.ReadFile(name); err != nil {
				r.Logger.Errorf("failed reading: %s", name)
				continue
			}
			tmpl := template.New("__main__").Funcs(Funcs).Option("missingkey=zero")
			if tmpl, err = tmpl.Parse(string(buf)); err != nil {
				r.Logger.Errorf("failed loading %s: %s", name, err.Error())
				continue
			}
			out := &bytes.Buffer{}
			if err = tmpl.Execute(out, map[string]interface{}{
				"Env": env,
			}); err != nil {
				r.Logger.Errorf("failed rendering %s: %s", name, err.Error())
				continue
			}
			content := out.Bytes()
			if !r.Unit.Raw {
				content = sanitizeLines(content)
			}
			if err = os.WriteFile(name, content, 0755); err != nil {
				r.Logger.Errorf("failed writing %s: %s", name, err.Error())
				continue
			}
			r.Logger.Printf("render finished: %s", name)
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
