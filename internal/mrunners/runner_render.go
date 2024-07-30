package mrunners

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/yankeguo/minit/internal/menv"
	"github.com/yankeguo/minit/internal/mtmpl"
	"github.com/yankeguo/minit/internal/munit"
	"github.com/yankeguo/rg"
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

func (r *actionRender) doFile(ctx context.Context, src string, dst string, env map[string]string) (err error) {

	var buf []byte
	if buf, err = os.ReadFile(src); err != nil {
		err = fmt.Errorf("failed reading %s: %s", src, err.Error())
		return
	}

	var content []byte
	if content, err = mtmpl.Execute(string(buf), map[string]any{
		"Env": env,
	}); err != nil {
		err = fmt.Errorf("failed rendering %s: %s", src, err.Error())
		return
	}

	if !r.Unit.Raw {
		content = sanitizeLines(content)
	}

	var dstPerm os.FileMode = 0644

	if src != dst {
		var srcFileInfo os.FileInfo
		if srcFileInfo, err = os.Stat(src); err != nil {
			err = fmt.Errorf("failed stating %s: %s", src, err.Error())
			return
		}
		dstPerm = srcFileInfo.Mode().Perm()

		srcDir := filepath.Dir(src)

		var srcDirInfo os.FileInfo
		if srcDirInfo, err = os.Stat(srcDir); err != nil {
			err = fmt.Errorf("failed stating %s: %s", srcDir, err.Error())
			return
		}

		if err = os.MkdirAll(filepath.Dir(dst), srcDirInfo.Mode().Perm()); err != nil {
			err = fmt.Errorf("failed mkdir %s: %s", filepath.Dir(dst), err.Error())
			return
		}
	}

	if err = os.WriteFile(dst, content, dstPerm); err != nil {
		err = fmt.Errorf("failed writing %s: %s", dst, err.Error())
		return
	}
	return
}

func (r *actionRender) Do(ctx context.Context) (err error) {
	r.Print("controller started")
	defer r.Print("controller exited")
	defer rg.Guard(&err)

	var env map[string]string

	if env, err = menv.Construct(menv.Environ(), r.Unit.Env); err != nil {
		err = r.PanicOnCritical("failed constructing environments variables", err)
		return
	}

	tasks := map[[2]string]struct{}{}

	for _, filePattern := range r.Unit.Files {

		segments := strings.Split(filePattern, ":")

		if len(segments) == 3 {
			var (
				dirSrc = strings.TrimSpace(segments[0])
				match  = strings.TrimSpace(segments[1])
				dirDst = strings.TrimSpace(segments[2])
			)
			if dirSrc == "" || match == "" || dirDst == "" {
				err = r.PanicOnCritical("failed parsing file pattern", errors.New("invalid file pattern: "+filePattern))
				continue
			}

			var names []string

			if names, err = filepath.Glob(filepath.Join(dirSrc, match)); err != nil {
				err = r.PanicOnCritical("failed globbing", err)
				continue
			}

			for _, name := range names {
				var relPath string
				if relPath, err = filepath.Rel(dirSrc, name); err != nil {
					err = r.PanicOnCritical("failed getting relative path", err)
					continue
				}
				tasks[[2]string{name, filepath.Join(dirDst, relPath)}] = struct{}{}
			}
		} else if len(segments) == 2 {
			var (
				fileSrc = strings.TrimSpace(segments[0])
				fileDst = strings.TrimSpace(segments[1])
			)
			if fileSrc == "" || fileDst == "" {
				err = r.PanicOnCritical("failed parsing file pattern", errors.New("invalid file pattern: "+filePattern))
				continue
			}
			tasks[[2]string{fileSrc, fileDst}] = struct{}{}
		} else if len(segments) == 1 {
			var names []string

			if names, err = filepath.Glob(filePattern); err != nil {
				err = r.PanicOnCritical("failed globbing", err)
				continue
			}

			for _, name := range names {
				tasks[[2]string{name, name}] = struct{}{}
			}
		} else {
			err = r.PanicOnCritical("failed parsing file pattern", errors.New("invalid file pattern: "+filePattern))
			continue
		}
	}

	for task := range tasks {
		if err = r.doFile(ctx, task[0], task[1], env); err != nil {
			err = r.PanicOnCritical("failed rendering", err)
			continue
		}
		r.Print("done rendering: " + task[1])
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
