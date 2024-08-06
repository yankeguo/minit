package mrunners

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
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
		defer rg.Guard(&err)
		rg.Must0(opts.Unit.RequireFiles())

		runner.Order = 10
		runner.Action = &actionRender{RunnerOptions: opts}
		return
	})
}

type actionRender struct {
	RunnerOptions
}

func (r *actionRender) doFile(_ context.Context, src string, dst string, env map[string]string) (err error) {
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

func (r *actionRender) addTask3(tasks map[[2]string]struct{}, dirSrc, match, dirDst string) (err error) {
	if dirSrc == "" || match == "" || dirDst == "" {
		err = r.PanicOnCritical("failed parsing file pattern", errors.New("invalid file pattern: "+fmt.Sprintf("%s:%s:%s", dirSrc, match, dirDst)))
		return
	}

	var names []string

	if names, err = filepath.Glob(filepath.Join(dirSrc, match)); err != nil {
		err = r.PanicOnCritical("failed globbing", err)
		return
	}

	for _, name := range names {
		var relPath string
		if relPath, err = filepath.Rel(dirSrc, name); err != nil {
			err = r.PanicOnCritical("failed getting relative path", err)
			continue
		}
		tasks[[2]string{name, filepath.Join(dirDst, relPath)}] = struct{}{}
	}
	return
}

func (r *actionRender) addTask2(tasks map[[2]string]struct{}, fileSrc, fileDst string) (err error) {
	var (
		infoSrc fs.FileInfo
		infoDst fs.FileInfo
	)

	if fileSrc == "" || fileDst == "" {
		err = r.PanicOnCritical("failed parsing file pattern", errors.New("invalid file pattern: "+fmt.Sprintf("%s:%s", fileSrc, fileDst)))
		return
	}

	if infoSrc, err = os.Stat(fileSrc); err != nil {
		err = r.PanicOnCritical("failed stating source file/directory", err)
		return
	}

	if infoDst, err = os.Stat(fileDst); err != nil {
		if os.IsNotExist(err) {
			err = nil
			infoDst = nil
		} else {
			err = r.PanicOnCritical("failed stating destination file", err)
			return
		}
	}

	if infoSrc.IsDir() && (infoDst == nil || infoDst.IsDir()) {
		err = r.addTask3(tasks, fileSrc, "*", fileDst)
		return
	} else if !infoSrc.IsDir() && (infoDst == nil || !infoDst.IsDir()) {
		tasks[[2]string{fileSrc, fileDst}] = struct{}{}
	} else {
		err = r.PanicOnCritical("failed checking source and destination", errors.New("source and destination must be both file or directory"))
		return
	}
	return
}

func (r *actionRender) addTask1(tasks map[[2]string]struct{}, filePattern string) (err error) {
	var names []string

	if names, err = filepath.Glob(filePattern); err != nil {
		err = r.PanicOnCritical("failed globbing", err)
		return
	}

	for _, name := range names {
		tasks[[2]string{name, name}] = struct{}{}
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
			if err = r.addTask3(
				tasks,
				strings.TrimSpace(segments[0]),
				strings.TrimSpace(segments[1]),
				strings.TrimSpace(segments[2]),
			); err != nil {
				err = r.PanicOnCritical("failed parsing file pattern", err)
				continue
			}
		} else if len(segments) == 2 {
			if err = r.addTask2(
				tasks,
				strings.TrimSpace(segments[0]),
				strings.TrimSpace(segments[1]),
			); err != nil {
				err = r.PanicOnCritical("failed parsing file pattern", err)
				continue
			}
		} else if len(segments) == 1 {
			if err = r.addTask1(tasks, filePattern); err != nil {
				err = r.PanicOnCritical("failed parsing file pattern", err)
				continue
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
