package mtmpl

import (
	"bytes"
	"text/template"
)

// Execute render text template with predefined funcs
func Execute(src string, data any) (out []byte, err error) {
	var t *template.Template
	if t, err = template.
		New("__main__").
		Funcs(Funcs).
		Option("missingkey=zero").
		Parse(src); err != nil {
		return
	}
	o := &bytes.Buffer{}
	if err = t.Execute(o, data); err != nil {
		return
	}
	out = o.Bytes()
	return
}
