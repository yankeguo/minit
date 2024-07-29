package munit

import "strings"

// FilterMap is a map of unit filters.
type FilterMap map[string]struct{}

// Blank returns true if the FilterMap is empty.
func (fm FilterMap) Blank() bool {
	return len(fm) == 0
}

// Match returns true if the unit matches the FilterMap.
func (fm FilterMap) Match(unit Unit) bool {
	if fm.Blank() {
		return false
	}
	if _, ok := fm[unit.Name]; ok {
		return true
	}
	if _, ok := fm[PrefixGroup+unit.Group]; ok {
		return true
	}
	if _, ok := fm[PrefixKind+unit.Kind]; ok {
		return true
	}
	return false
}

// NewFilterMap creates a new FilterMap from a comma separated string.
func NewFilterMap(s string) (out FilterMap) {
	out = FilterMap{}
	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if item == "" || item == PrefixGroup || item == PrefixKind {
			continue
		}
		out[item] = struct{}{}
	}
	return
}

// Filter is a filter for units, it has either a pass or a deny filter map.
type Filter struct {
	pass FilterMap
	deny FilterMap
}

// NewFilter creates a new Filter from either a pass or a deny string.
func NewFilter(pass, deny string) (uf *Filter) {
	return &Filter{
		pass: NewFilterMap(pass),
		deny: NewFilterMap(deny),
	}
}

// Match returns true if the unit matches the filter.
func (uf *Filter) Match(unit Unit) bool {
	if !uf.pass.Blank() {
		if !uf.pass.Match(unit) {
			return false
		}
	}
	if !uf.deny.Blank() {
		if uf.deny.Match(unit) {
			return false
		}
	}
	return true
}
