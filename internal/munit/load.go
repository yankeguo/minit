package munit

import (
	"fmt"
	"os"
	"regexp"
	"sort"
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
	Dirs []string
}

func Load(opts LoadOptions) (output []Unit, skipped []Unit, err error) {
	defer rg.Guard(&err)

	// create a filter
	filter := NewFilter("", "")

	if opts.Env != nil {
		filter = NewFilter(
			opts.Env["MINIT_ENABLE"],
			opts.Env["MINIT_DISABLE"],
		)
	}

	// load units in order of dirs, env, args
	var units []Unit

	for _, dir := range opts.Dirs {
		units = append(units, rg.Must(LoadDir(dir))...)
	}

	if opts.Env != nil {
		if unit, ok := rg.Must2(LoadEnv(opts.Env)); ok {
			units = append(units, unit)
		}

		for _, infix := range DetectEnvInfixes(opts.Env) {
			if unit, ok := rg.Must2(LoadEnvWithInfix(opts.Env, infix)); ok {
				units = append(units, unit)
			}
		}
	}

	if len(opts.Args) > 0 {
		if unit, ok := rg.Must2(LoadArgs(opts.Args)); ok {
			units = append(units, unit)
		}
	}

	sortUnits(units)

	// check duplicated name
	names := map[string]struct{}{}

	// reserve 'minit'

	names[NameMinit] = struct{}{}

	// whitelist / blacklist, replicas
	for _, unit := range units {
		// check unit kind
		if _, ok := knownUnitKind[unit.Kind]; !ok {
			err = fmt.Errorf("invalid unit kind '%s' for unit '%s': must be one of: render, once, daemon, cron", unit.Kind, unit.Name)
			return
		}

		// check unit name
		if !regexpName.MatchString(unit.Name) {
			err = fmt.Errorf("invalid unit name '%s': name must start with a letter, contain only alphanumeric characters, hyphens, or underscores, and end with an alphanumeric character", unit.Name)
			return
		}

		// check duplicated
		if _, found := names[unit.Name]; found {
			err = fmt.Errorf("duplicated unit name '%s': each unit must have a unique name", unit.Name)
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

func sortUnits(units []Unit) {
	// render first
	sort.SliceStable(units, func(i, j int) bool {
		return units[i].Kind == KindRender && units[j].Kind != KindRender
	})
	// then by order
	sort.SliceStable(units, func(i, j int) bool {
		return units[i].Order < units[j].Order
	})
}
