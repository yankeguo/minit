package mtmpl

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

func doTestTemplate(t *testing.T, want string, src string, data interface{}) {
	tpl, err := template.New("__main__").Funcs(Funcs).Option("missingkey=zero").Parse(src)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, data)
	require.NoError(t, err)
	require.Equal(t, want, strings.TrimSpace(buf.String()))
}

func TestFuncNetResolveIP(t *testing.T) {
	doTestTemplate(t, `127.0.0.1`, `{{netResolveIP "127.0.0.1"}}`, nil)
}

func TestFuncAdd(t *testing.T) {
	doTestTemplate(t, `3`, `{{add 1 2}}`, nil)
	doTestTemplate(t, `true`, `{{add true false}}`, nil)
}

func TestFuncNeg(t *testing.T) {
	doTestTemplate(t, `-1`, `{{neg 1}}`, nil)
	doTestTemplate(t, `false`, `{{neg true}}`, nil)
}

func TestFuncDic(t *testing.T) {
	doTestTemplate(t, `bar`, `
	{{define "hello"}}
	{{.foo}}
	{{end}}
	{{template "hello" (dict "foo" "bar")}}
	`, nil)
}

func TestFuncSlice(t *testing.T) {
	doTestTemplate(t, `1 2 3`, `{{range $i, $v := slice 1 2 3}}{{$v}} {{end}}`, nil)
	doTestTemplate(t, `3`, `
	{{index (slice 1 2 3) (add 1 1)}}
	`, nil)
}

func TestFuncFloat64(t *testing.T) {
	doTestTemplate(t, `4`, `{{add (float64 (int64 3)) (float64 1)}}`, nil)
}

func TestFuncReadFile(t *testing.T) {
	doTestTemplate(
		t,
		`worldb`,
		`{{ osReadFileString (filepathJoin "testdata" "hello.txt") | stringsTrimSpace }}b`,
		nil,
	)
}
