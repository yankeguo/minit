package munit

import (
	"errors"
	"strconv"
	"strings"

	"github.com/yankeguo/minit/pkg/shellquote"
)

const (
	EnvPrefixUnit = "MINIT_UNIT_"
)

// DetectEnvInfixes detects infixes from environment variables
func DetectEnvInfixes() (infixes []string) {
	_infixes := map[string]struct{}{}

	for _, item := range osEnviron() {
		splits := strings.SplitN(item, "=", 2)
		if len(splits) != 2 {
			continue
		}

		key := splits[0]
		if !strings.HasPrefix(key, EnvPrefixUnit) {
			continue
		}
		key = strings.TrimPrefix(key, EnvPrefixUnit)
		if strings.HasSuffix(key, "_COMMAND") {
			// anything with _COMMAND is a unit
			key = strings.TrimSuffix(key, "_COMMAND")
			if key == "" {
				continue
			}
			_infixes[key] = struct{}{}
		} else if strings.HasSuffix(key, "_FILES") {
			// anything with _FILES and has _KIND set to render is a render unit
			key = strings.TrimSuffix(key, "_FILES")
			if key == "" {
				continue
			}
			if osGetenv(EnvPrefixUnit+key+"_KIND") == KindRender {
				_infixes[key] = struct{}{}
			}
		}
	}

	for infix := range _infixes {
		infixes = append(infixes, infix)
	}

	return
}

// LoadEnvWithInfix loads unit from environment variables with infix
// e.g. MINIT_UNIT_HELLO_KIND, MINIT_UNIT_HELLO_NAME, MINIT_UNIT_HELLO_COMMAND
func LoadEnvWithInfix(infix string) (unit Unit, ok bool, err error) {
	// kind
	unit.Kind = osGetenv(EnvPrefixUnit + infix + "_KIND")

	if unit.Kind == "" {
		unit.Kind = KindDaemon
	}

	switch unit.Kind {
	case KindDaemon, KindOnce, KindCron, KindRender:
	default:
		err = errors.New("unsupported $" + EnvPrefixUnit + infix + "_KIND: " + unit.Kind)
		return
	}

	// name, group, count
	unit.Name = osGetenv(EnvPrefixUnit + infix + "_NAME")
	if unit.Name == "" {
		unit.Name = "env-" + strings.ToLower(infix)
	}
	unit.Group = osGetenv(EnvPrefixUnit + infix + "_GROUP")
	unit.Count, _ = strconv.Atoi(osGetenv(EnvPrefixUnit + infix + "_COUNT"))

	// command, dir, shell, charset
	switch unit.Kind {
	case KindDaemon, KindOnce, KindCron:
		cmd := osGetenv(EnvPrefixUnit + infix + "_COMMAND")

		if unit.Command, err = shellquote.Split(cmd); err != nil {
			return
		}

		if len(unit.Command) == 0 {
			err = errors.New("missing environment variable $MINIT_" + infix + "_COMMAND")
			return
		}

		unit.Dir = osGetenv(EnvPrefixUnit + infix + "_DIR")
		unit.Shell = osGetenv(EnvPrefixUnit + infix + "_SHELL")
		unit.Charset = osGetenv(EnvPrefixUnit + infix + "_CHARSET")

		for _, item := range strings.Split(osGetenv(EnvPrefixUnit+infix+"_ENV"), ";") {
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
		unit.Cron = osGetenv(EnvPrefixUnit + infix + "_CRON")

		if unit.Cron == "" {
			err = errors.New("missing environment variable $" + EnvPrefixUnit + infix + "_CRON while $" + EnvPrefixUnit + infix + "_KIND is 'cron'")
			return
		}

		unit.Immediate, _ = strconv.ParseBool(osGetenv(EnvPrefixUnit + infix + "_IMMEDIATE"))
	}

	// raw, files
	if unit.Kind == KindRender {
		unit.Raw, _ = strconv.ParseBool(osGetenv(EnvPrefixUnit + infix + "_RAW"))

		for _, item := range strings.Split(osGetenv(EnvPrefixUnit+infix+"_FILES"), ";") {
			item = strings.TrimSpace(item)
			if item != "" {
				unit.Files = append(unit.Files, item)
			}
		}

		if len(unit.Files) == 0 {
			err = errors.New("missing environment variable $" + EnvPrefixUnit + infix + "_FILES while $" + EnvPrefixUnit + infix + "_KIND is 'render'")
			return
		}
	}

	// blocking
	if unit.Kind == KindOnce {
		if nb, err := strconv.ParseBool(strings.TrimSpace(osGetenv(EnvPrefixUnit + infix + "_BLOCKING"))); err == nil && !nb {
			unit.Blocking = new(bool)
			*unit.Blocking = false
		}
	}

	// critical
	unit.Critical, _ = strconv.ParseBool(osGetenv(EnvPrefixUnit + infix + "_CRITICAL"))

	// success codes
	for _, item := range strings.Split(osGetenv(EnvPrefixUnit+infix+"_SUCCESS_CODES"), ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if code, err := strconv.Atoi(item); err == nil {
			unit.SuccessCodes = append(unit.SuccessCodes, code)
		}
	}

	ok = true

	return
}

// LoadEnv loads legacy main unit from environment variables
func LoadEnv() (unit Unit, ok bool, err error) {
	cmd := strings.TrimSpace(osGetenv("MINIT_MAIN"))
	if cmd == "" {
		return
	}

	name := strings.TrimSpace(osGetenv("MINIT_MAIN_NAME"))
	if name == "" {
		name = "env-main"
	}

	var (
		cron      string
		immediate bool
	)

	kind := strings.TrimSpace(osGetenv("MINIT_MAIN_KIND"))

	switch kind {
	case KindDaemon, KindOnce:
	case KindCron:
		cron = strings.TrimSpace(osGetenv("MINIT_MAIN_CRON"))

		if cron == "" {
			err = errors.New("missing environment variable $MINIT_MAIN_CRON while $MINIT_MAIN_KIND is 'cron'")
			return
		}

		immediate, _ = strconv.ParseBool(osGetenv("MINIT_MAIN_IMMEDIATE"))
	case "":
		if once, _ := strconv.ParseBool(strings.TrimSpace(osGetenv("MINIT_MAIN_ONCE"))); once {
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
		Group:     strings.TrimSpace(osGetenv("MINIT_MAIN_GROUP")),
		Kind:      kind,
		Cron:      cron,
		Immediate: immediate,
		Command:   cmds,
		Dir:       strings.TrimSpace(osGetenv("MINIT_MAIN_DIR")),
		Charset:   strings.TrimSpace(osGetenv("MINIT_MAIN_CHARSET")),
	}

	if unit.Kind == KindOnce {
		if nb, err := strconv.ParseBool(strings.TrimSpace(osGetenv("MINIT_MAIN_BLOCKING"))); err == nil && !nb {
			unit.Blocking = new(bool)
			*unit.Blocking = false
		}
	}

	ok = true
	return
}