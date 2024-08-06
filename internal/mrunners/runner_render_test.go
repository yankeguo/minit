package mrunners

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/minit/internal/mexec"
	"github.com/yankeguo/minit/internal/mlog"
	"github.com/yankeguo/minit/internal/munit"
	"github.com/yankeguo/rg"
)

func TestRunnerRender(t *testing.T) {
	dirSrc := t.TempDir()
	os.MkdirAll(filepath.Join(dirSrc, "dir2"), 0755)
	os.WriteFile(filepath.Join(dirSrc, "file11.txt"), []byte("{{stringsToUpper .Env.Foo}}"), 0755)
	os.WriteFile(filepath.Join(dirSrc, "file12.txt"), []byte("{{stringsToUpper .Env.Foo}}"), 0755)
	os.WriteFile(filepath.Join(dirSrc, "file21.txt"), []byte("{{stringsToUpper .Env.Bar}}"), 0755)
	os.WriteFile(filepath.Join(dirSrc, "dir2", "file22.txt"), []byte("{{stringsToUpper .Env.Bar}}"), 0755)
	os.WriteFile(filepath.Join(dirSrc, "file31.txt"), []byte("{{stringsToUpper .Env.Foo}}  \n          {{stringsToUpper .Env.Bar}}\n\n  \n{{stringsToUpper .Env.Foo}}\n"), 0755)
	os.WriteFile(filepath.Join(dirSrc, "file32.txt"), []byte("{{stringsToUpper .Env.Foo}}  \n          {{stringsToUpper .Env.Bar}}\n\n  \n{{stringsToUpper .Env.Foo}}\n"), 0755)

	dirDst := t.TempDir()

	exem := mexec.NewManager()

	r := &actionRender{
		RunnerOptions: RunnerOptions{
			Unit: munit.Unit{
				Kind: munit.KindRender,
				Name: "test",
				Files: []string{
					filepath.Join(dirSrc, "file1*.txt"),
					filepath.Join(dirSrc, "file21.txt") + ":" + filepath.Join(dirDst, "file211.txt"),
					filepath.Join(dirSrc, "dir2") + ":" + filepath.Join(dirDst, "dir2"),
					dirSrc + ":file3*:" + dirDst,
				},
				Critical: true,
				Env: map[string]string{
					"Foo": "foo",
					"Bar": "bar",
				},
			},
			Exec:   exem,
			Logger: rg.Must(mlog.NewProcLogger(mlog.ProcLoggerOptions{})),
		},
	}
	err := r.Do(context.Background())
	require.NoError(t, err)
	buf1, err := os.ReadFile(filepath.Join(dirSrc, "file11.txt"))
	require.NoError(t, err)
	require.Equal(t, "FOO", string(buf1))
	buf2, err := os.ReadFile(filepath.Join(dirDst, "dir2", "file22.txt"))
	require.NoError(t, err)
	require.Equal(t, "BAR", string(buf2))
	buf3, err := os.ReadFile(filepath.Join(dirDst, "file32.txt"))
	require.NoError(t, err)
	require.Equal(t, "FOO\n          BAR\nFOO", string(buf3))
}
