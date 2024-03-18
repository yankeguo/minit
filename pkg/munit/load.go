package munit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yankeguo/minit/pkg/shellquote"
	"gopkg.in/yaml.v3"
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

const (
	UnitPrefix = "MINIT_UNIT_"
)

// DetectEnvInfixes detects infixes from environment variables
func DetectEnvInfixes() (infixes []string) {
	_infixes := map[string]struct{}{}

	for _, item := range os.Environ() {
		splits := strings.SplitN(item, "=", 2)
		if len(splits) != 2 {
			continue
		}

		key := splits[0]
		if !strings.HasPrefix(key, UnitPrefix) {
			continue
		}
		key = strings.TrimPrefix(key, UnitPrefix)
		if strings.HasSuffix(key, "_COMMAND") {
			key = strings.TrimSuffix(key, "_COMMAND")
			if key == "" {
				continue
			}
			_infixes[key] = struct{}{}
		} else if strings.HasSuffix(key, "_FILES") {
			key = strings.TrimSuffix(key, "_FILES")
			if key == "" {
				continue
			}
			if os.Getenv(UnitPrefix+key+"_KIND") == KindRender {
				_infixes[key] = struct{}{}
			}
		}
	}

	for infix := range _infixes {
		infixes = append(infixes, infix)
	}

	return
}

// LoadFromEnvWithInfix loads unit from environment variables with infix
// e.g. MINIT_MAIN_KIND, MINIT_MAIN_NAME, MINIT_MAIN_COMMAND
func LoadFromEnvWithInfix(infix string) (unit Unit, ok bool, err error) {
	// kind
	unit.Kind = os.Getenv(UnitPrefix + infix + "_KIND")

	if unit.Kind == "" {
		unit.Kind = KindDaemon
	}

	switch unit.Kind {
	case KindDaemon, KindOnce, KindCron, KindRender:
	default:
		err = errors.New("unsupported $" + UnitPrefix + infix + "_KIND: " + unit.Kind)
		return
	}

	// name, group, count
	unit.Name = os.Getenv(UnitPrefix + infix + "_NAME")
	if unit.Name == "" {
		unit.Name = "env-" + strings.ToLower(infix)
	}
	unit.Group = os.Getenv(UnitPrefix + infix + "_GROUP")
	unit.Count, _ = strconv.Atoi(os.Getenv("MINIT_" + infix + "_COUNT"))

	// command, dir, shell, charset
	switch unit.Kind {
	case KindDaemon, KindOnce, KindCron:
		cmd := os.Getenv(UnitPrefix + infix + "_COMMAND")

		if unit.Command, err = shellquote.Split(cmd); err != nil {
			return
		}

		if len(unit.Command) == 0 {
			err = errors.New("missing environment variable $MINIT_" + infix + "_COMMAND")
			return
		}

		unit.Dir = os.Getenv(UnitPrefix + infix + "_DIR")
		unit.Shell = os.Getenv(UnitPrefix + infix + "_SHELL")
		unit.Charset = os.Getenv(UnitPrefix + infix + "_CHARSET")

		for _, item := range strings.Split(os.Getenv(UnitPrefix+infix+"_ENV"), ";") {
			item = strings.TrimSpace(item)
			splits := strings.SplitN(item, "=", 2)
			if len(splits) == 2 {
				if unit.Env == nil {
					unit.Env = make(map[string]string)
				}
				unit.Env[splits[0]] = splits[1]
			}
		}
	}

	// cron, immediate
	if unit.Kind == KindCron {
		unit.Cron = os.Getenv(UnitPrefix + infix + "_CRON")

		if unit.Cron == "" {
			err = errors.New("missing environment variable $" + UnitPrefix + infix + "_CRON while $" + UnitPrefix + infix + "_KIND is 'cron'")
			return
		}

		unit.Immediate, _ = strconv.ParseBool(os.Getenv(UnitPrefix + infix + "_IMMEDIATE"))
	}

	// raw, files
	if unit.Kind == KindRender {
		unit.Raw, _ = strconv.ParseBool(os.Getenv(UnitPrefix + infix + "_RAW"))

		for _, item := range strings.Split(os.Getenv(UnitPrefix+infix+"_FILES"), ";") {
			item = strings.TrimSpace(item)
			if item != "" {
				unit.Files = append(unit.Files, item)
			}
		}

		if len(unit.Files) == 0 {
			err = errors.New("missing environment variable $" + UnitPrefix + infix + "_FILES while $" + UnitPrefix + infix + "_KIND is 'render'")
			return
		}
	}

	// blocking
	if unit.Kind == KindOnce {
		if nb, err := strconv.ParseBool(strings.TrimSpace(os.Getenv(UnitPrefix + infix + "_BLOCKING"))); err == nil && !nb {
			unit.Blocking = new(bool)
			*unit.Blocking = false
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

	var (
		cron      string
		immediate bool
	)

	kind := strings.TrimSpace(os.Getenv("MINIT_MAIN_KIND"))

	switch kind {
	case KindDaemon, KindOnce:
	case KindCron:
		cron = strings.TrimSpace(os.Getenv("MINIT_MAIN_CRON"))

		if cron == "" {
			err = errors.New("missing environment variable $MINIT_MAIN_CRON while $MINIT_MAIN_KIND is 'cron'")
			return
		}

		immediate, _ = strconv.ParseBool(os.Getenv("MINIT_MAIN_IMMEDIATE"))
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
		Name:      name,
		Group:     strings.TrimSpace(os.Getenv("MINIT_MAIN_GROUP")),
		Kind:      kind,
		Cron:      cron,
		Immediate: immediate,
		Command:   cmds,
		Dir:       strings.TrimSpace(os.Getenv("MINIT_MAIN_DIR")),
		Charset:   strings.TrimSpace(os.Getenv("MINIT_MAIN_CHARSET")),
	}

	if unit.Kind == KindOnce {
		if nb, err := strconv.ParseBool(strings.TrimSpace(os.Getenv("MINIT_MAIN_BLOCKING"))); err == nil && !nb {
			unit.Blocking = new(bool)
			*unit.Blocking = false
		}
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
