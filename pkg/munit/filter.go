package munit

import "strings"

type FilterMap map[string]struct{}

func (fm FilterMap) Match(unit Unit) bool {
	if fm == nil {
		return false
	}
	if _, ok := fm[unit.Name]; ok {
		return true
	}
	if _, ok := fm[PrefixGroup+unit.Group]; ok {
		return true
	}
	return false
}

func NewFilterMap(s string) (out FilterMap) {
	s = strings.TrimSpace(s)
	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if item == "" || item == PrefixGroup {
			continue
		}
		if out == nil {
			out = FilterMap{}
		}
		out[item] = struct{}{}
	}
	return
}

type Filter struct {
	pass FilterMap
	deny FilterMap
}

func NewFilter(pass, deny string) (uf *Filter) {
	return &Filter{
		pass: NewFilterMap(pass),
		deny: NewFilterMap(deny),
	}
}

func (uf *Filter) Match(unit Unit) bool {
	if uf.pass != nil {
		if !uf.pass.Match(unit) {
			return false
		}
	}
	if uf.deny != nil {
		if uf.deny.Match(unit) {
			return false
		}
	}
	return true
}
