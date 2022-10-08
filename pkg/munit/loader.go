package munit

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	regexpName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9]$`)
)

const (
	NameMinit = "minit"

	PrefixGroup = "@"
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
		var unit Unit
		var ok bool
		if unit, ok, err = LoadEnv(); err != nil {
			return
		}
		if ok {
			units = append(units, unit)
		}
	}
	if opts.Dir != "" {
		var dUnits []Unit
		if dUnits, err = LoadDir(opts.Dir); err != nil {
			return
		}
		units = append(units, dUnits...)
	}

	// check duplicated name
	names := map[string]struct{}{}

	// reserve 'minit'

	names[NameMinit] = struct{}{}

	// whitelist / blacklist, replicas
	for _, unit := range units {
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

		// replicas
		if unit.Count > 1 {
			for i := 0; i < unit.Count; i++ {
				subUnit := unit
				subUnit.Name = fmt.Sprintf("%s-%d", unit.Name, i+1)
				subUnit.Count = 1
				output = append(output, subUnit)
			}
		} else {
			unit.Count = 1
			output = append(output, unit)
		}
	}

	return
}
