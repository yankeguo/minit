package munit

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexpName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9]$`)
)

const (
	NameMinit = "minit"

	PrefixGroup = "@"
	PrefixKind  = "&"
)

type Loader struct {
	filter *Filter
}

func NewLoader() (ld *Loader) {
	return &Loader{
		filter: NewFilter(
			strings.TrimSpace(os.Getenv("MINIT_ENABLE")),
			strings.TrimSpace(os.Getenv("MINIT_DISABLE")),
		),
	}
}

type LoadOptions struct {
	Args []string
	Env  bool
	Dir  string
}

func (ld *Loader) Load(opts LoadOptions) (output []Unit, skipped []Unit, err error) {
	var units []Unit

	// load units
	if opts.Dir != "" {
		var dUnits []Unit
		if dUnits, err = LoadDir(opts.Dir); err != nil {
			return
		}
		units = append(units, dUnits...)
	}
	if len(opts.Args) > 0 {
		var unit Unit
		var ok bool
		if unit, ok, err = LoadArgs(opts.Args); err != nil {
			return
		}
		if ok {
			units = append(units, unit)
		}
	}
	if opts.Env {
		{
			// legacy minit main
			var (
				unit Unit
				ok   bool
			)
			if unit, ok, err = LoadEnv(); err != nil {
				return
			} else if ok {
				units = append(units, unit)
			}
		}

		for _, infix := range DetectEnvInfixes() {
			var (
				unit Unit
				ok   bool
			)
			if unit, ok, err = LoadFromEnvWithInfix(infix); err != nil {
				return
			}
			if ok {
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
		if !ld.filter.Match(unit) {
			skipped = append(skipped, unit)
			continue
		}

		// eval cron
		if unit.Cron != "" {
			unit.Cron = os.ExpandEnv(unit.Cron)
		}

		// replicas
		if unit.Count > 1 {
			for i := 0; i < unit.Count; i++ {
				subUnit := unit
				subUnit.Name = unit.Name + "-" + strconv.Itoa(i+1)
				subUnit.Count = 1
				dupOrMakeMap(&subUnit.Env)
				subUnit.Env["MINIT_UNIT_NAME"] = subUnit.Name
				subUnit.Env["MINIT_UNIT_SUB_ID"] = strconv.Itoa(i + 1)

				output = append(output, subUnit)
			}
		} else {
			unit.Count = 1
			dupOrMakeMap(&unit.Env)
			unit.Env["MINIT_UNIT_NAME"] = unit.Name
			unit.Env["MINIT_UNIT_SUB_ID"] = "1"

			output = append(output, unit)
		}
	}

	return
}

func dupOrMakeMap[T comparable, U any](m *map[T]U) {
	nm := make(map[T]U)
	if *m != nil {
		for k, v := range *m {
			nm[k] = v
		}
	}
	*m = nm
}
