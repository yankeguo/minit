package munit

import (
	"errors"
	"fmt"
	"github.com/guoyk93/minit/pkg/shellquote"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func LoadArgs(args []string) (unit Unit, ok bool, err error) {
	var opts []string

	// fix a history issue
	for len(args) > 0 {
		if filepath.Base(args[0]) == "minit" {
			args = args[1:]
		} else {
			break
		}
	}

	// extract arguments after '--' if existed
	for i, item := range args {
		if item == "--" {
			opts = args[0:i]
			args = args[i+1:]
			break
		}
	}

	if len(args) == 0 {
		return
	}

	unit = Unit{
		Name:    "arg-main",
		Kind:    KindDaemon,
		Command: args,
	}

	// opts decoding
	for _, opt := range opts {
		if strings.HasSuffix(opt, "-"+KindOnce) {
			unit.Kind = KindOnce
		}
	}

	ok = true

	return
}

func LoadEnv() (unit Unit, ok bool, err error) {
	cmd := strings.TrimSpace(os.Getenv("MINIT_MAIN"))
	if cmd == "" {
		return
	}

	name := strings.TrimSpace(os.Getenv("MINIT_MAIN_NAME"))
	if name == "" {
		name = "env-main"
	}

	var cron string

	kind := strings.TrimSpace(os.Getenv("MINIT_MAIN_KIND"))

	switch kind {
	case KindDaemon, KindOnce:
	case KindCron:
		cron = strings.TrimSpace(os.Getenv("MINIT_MAIN_CRON"))

		if cron == "" {
			err = errors.New("missing environment variable $MINIT_MAIN_CRON while $MINIT_MAIN_KIND is 'cron'")
			return
		}
	case "":
		if once, _ := strconv.ParseBool(strings.TrimSpace(os.Getenv("MINIT_MAIN_ONCE"))); once {
			kind = KindOnce
		} else {
			kind = KindDaemon
		}
	default:
		err = errors.New("unsupported $MINIT_MAIN_KIND: " + kind)
		return
	}

	var cmds []string
	if cmds, err = shellquote.Split(cmd); err != nil {
		return
	}

	unit = Unit{
		Name:    name,
		Group:   strings.TrimSpace(os.Getenv("MINIT_MAIN_GROUP")),
		Kind:    kind,
		Cron:    cron,
		Command: cmds,
		Dir:     strings.TrimSpace(os.Getenv("MINIT_MAIN_DIR")),
		Charset: strings.TrimSpace(os.Getenv("MINIT_MAIN_CHARSET")),
	}

	ok = true
	return
}

func LoadFile(filename string) (units []Unit, err error) {
	var f *os.File
	if f, err = os.Open(filename); err != nil {
		return
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	for {
		var unit Unit
		if err = dec.Decode(&unit); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = fmt.Errorf("failed to decode unit file %s: %s", filename, err.Error())
			}
			return
		}

		if unit.Kind == "" {
			continue
		}

		units = append(units, unit)
	}
}

func LoadDir(dir string) (units []Unit, err error) {
	for _, ext := range []string{"*.yml", "*.yaml"} {
		var files []string
		if files, err = filepath.Glob(filepath.Join(dir, ext)); err != nil {
			return
		}
		for _, file := range files {
			var _units []Unit
			if _units, err = LoadFile(file); err != nil {
				return
			}
			units = append(units, _units...)
		}
	}
	return
}
