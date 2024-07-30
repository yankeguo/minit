package mrunners

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/minit/pkg/mexec"
	"github.com/yankeguo/minit/pkg/mlog"
	"github.com/yankeguo/minit/pkg/munit"
	"github.com/yankeguo/rg"
)

func TestRunnerRender(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("{{stringsToUpper .Env.Foo}}"), 0755)
	os.WriteFile(filepath.Join(dir, "file2.txt"), []byte("{{stringsToUpper .Env.Bar}}"), 0755)
	os.WriteFile(filepath.Join(dir, "file3.txt"), []byte("{{stringsToUpper .Env.Foo}}  \n          {{stringsToUpper .Env.Bar}}\n\n  \n{{stringsToUpper .Env.Foo}}\n"), 0755)

	exem := mexec.NewManager()

	buf := &bytes.Buffer{}

	r := &runnerRender{
		RunnerOptions: RunnerOptions{
			Unit: munit.Unit{
				Kind:     munit.KindRender,
				Name:     "test",
				Files:    []string{filepath.Join(dir, "*.txt")},
				Critical: true,
				Env: map[string]string{
					"Foo": "foo",
					"Bar": "bar",
				},
			},
			Exec: exem,
			Logger: rg.Must(mlog.NewProcLogger(mlog.ProcLoggerOptions{
				ConsoleOut: buf,
				ConsoleErr: buf,
			})),
		},
	}
	err := r.Do(context.Background())
	require.NoError(t, err)
	buf1, err := os.ReadFile(filepath.Join(dir, "file1.txt"))
	require.NoError(t, err)
	require.Equal(t, "FOO", string(buf1))
	buf2, err := os.ReadFile(filepath.Join(dir, "file2.txt"))
	require.NoError(t, err)
	require.Equal(t, "BAR", string(buf2))
	buf3, err := os.ReadFile(filepath.Join(dir, "file3.txt"))
	require.NoError(t, err)
	require.Equal(t, "FOO\n          BAR\nFOO", string(buf3))
}
