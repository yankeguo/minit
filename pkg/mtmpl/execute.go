package mtmpl

import (
	"bytes"
	"text/template"
)

func Execute(src string, data any) (buf []byte, err error) {
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
	buf = o.Bytes()
	return
}
