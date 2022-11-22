package mtmpl

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"text/template"
)

const TEST_TMPL = `
{{$a := 3}}
{{$b := 1}}
{{add (neg $b) $a}}
{{.A}}
`

func TestFuncs(t *testing.T) {
	tmpl := template.New("__main__").Funcs(Funcs).Option("missingkey=zero")
	tmpl, err := tmpl.Parse(TEST_TMPL)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, map[string]interface{}{"A": "B"})
	require.NoError(t, err)
	require.Equal(t, "2\nB", strings.TrimSpace(buf.String()))
}
