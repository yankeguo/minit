package munit

import (
	"errors"
	"os"
	"regexp"
	"strconv"

	"github.com/yankeguo/rg"
)

var (
	regexpName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9]$`)
)

const (
	NameMinit = "minit"

	FilterPrefixGroup = "@"
	FilterPrefixKind  = "&"
)

type LoadOptions struct {
	Args []string
	Env  map[string]string
	Dir  string
}

func Load(opts LoadOptions) (output []Unit, skipped []Unit, err error) {
	defer rg.Guard(&err)

	var units []Unit

	if opts.Dir != "" {
		units = append(units, rg.Must(LoadDir(opts.Dir))...)
	}

	if len(opts.Args) > 0 {
		if unit, ok := rg.Must2(LoadArgs(opts.Args)); ok {
			units = append(units, unit)
		}
	}

	filter := NewFilter("", "")

	if opts.Env != nil {

		filter = NewFilter(
			opts.Env["MINIT_ENABLE"],
			opts.Env["MINIT_DISABLE"],
		)

		if unit, ok := rg.Must2(LoadEnv(opts.Env)); ok {
			units = append(units, unit)
		}

		for _, infix := range DetectEnvInfixes(opts.Env) {
			if unit, ok := rg.Must2(LoadEnvWithInfix(opts.Env, infix)); ok {
				units = append(units, unit)
			}
		}
	}

	// check duplicated name
	names := map[string]struct{}{}

	// reserve 'minit'

	names[NameMinit] = struct{}{}

	// whitelist / blacklist, replicas
	for _, unit := range units {
		// check unit kind
		if _, ok := knownUnitKind[unit.Kind]; !ok {
			err = errors.New("invalid unit kind: " + unit.Kind)
			return
		}

		// check unit name
		if !regexpName.MatchString(unit.Name) {
			err = errors.New("invalid unit name: " + unit.Name)
			return
		}

		// check duplicated
		if _, found := names[unit.Name]; found {
			err = errors.New("duplicated unit name: " + unit.Name)
			return
		}

		names[unit.Name] = struct{}{}

		// fix default group
		if unit.Group == "" {
			unit.Group = DefaultGroup
		}

		// skip if needed
		if !filter.Match(unit) {
			skipped = append(skipped, unit)
			continue
		}

		// eval cron
		if unit.Cron != "" && opts.Env != nil {
			unit.Cron = os.Expand(unit.Cron, func(s string) string {
				return opts.Env[s]
			})
		}

		// replicas
		if unit.Count > 1 {
			for i := 0; i < unit.Count; i++ {
				subUnit := unit
				subUnit.Name = unit.Name + "-" + strconv.Itoa(i+1)
				subUnit.Count = 1
				duplicateMap(&subUnit.Env)
				subUnit.Env["MINIT_UNIT_NAME"] = subUnit.Name
				subUnit.Env["MINIT_UNIT_SUB_ID"] = strconv.Itoa(i + 1)

				output = append(output, subUnit)
			}
		} else {
			unit.Count = 1
			duplicateMap(&unit.Env)
			unit.Env["MINIT_UNIT_NAME"] = unit.Name
			unit.Env["MINIT_UNIT_SUB_ID"] = "1"

			output = append(output, unit)
		}
	}

	return
}

func duplicateMap[T comparable, U any](m *map[T]U) {
	nm := make(map[T]U)
	if *m != nil {
		for k, v := range *m {
			nm[k] = v
		}
	}
	*m = nm
}
