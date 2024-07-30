package mlog

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestLog(t *testing.T) {
	os.MkdirAll(filepath.Join("testdata", "logger"), 0755)
	os.WriteFile(filepath.Join("testdata", "logger", ".gitignore"), []byte("*.log"), 0644)
	log, err := NewProcLogger(ProcLoggerOptions{
		FileOptions: &RotatingFileOptions{
			Dir:      filepath.Join("testdata", "logger"),
			Filename: "test",
		},
		ConsolePrefix: "test",
	})
	require.NoError(t, err)
	log.Print("hello", "world")
	log.Printf("hello, %s", "world")
	log.Error("error", "world")
	log.Errorf("error, %s", "world")
}
