package mexec

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/minit/internal/mlog"
)

func TestNewManager(t *testing.T) {
	m := NewManager()

	os.RemoveAll(filepath.Join("testdata", "test.out.log"))
	os.RemoveAll(filepath.Join("testdata", "test.err.log"))

	logger, err := mlog.NewProcLogger(mlog.ProcLoggerOptions{
		FileOptions: &mlog.RotatingFileOptions{
			Dir:      "testdata",
			Filename: "test",
		},
	})
	require.NoError(t, err)

	err = m.Execute(ExecuteOptions{
		Dir: "testdata",
		Env: map[string]string{
			"AAA": "BBB",
		},
		Command: []string{
			"echo", "$AAA",
		},
		Logger: logger,
	})
	require.NoError(t, err)

	buf, err := os.ReadFile(filepath.Join("testdata", "test.out.log"))
	require.Contains(t, string(buf), "BBB")
	require.NoError(t, err)

	go func() {
		time.Sleep(time.Second)
		m.Signal(syscall.SIGINT)
	}()

	t1 := time.Now()

	err = m.Execute(ExecuteOptions{
		Dir: "testdata",
		Env: map[string]string{
			"AAA": "10",
		},
		Command: []string{
			"sleep", "$AAA",
		},
		Logger:       logger,
		SuccessCodes: []int{0, 130, -1},
	})
	require.NoError(t, err)
	require.True(t, time.Since(t1) < time.Second*2)
}
