package mlog

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestNewRotatingFile(t *testing.T) {
	_ = os.RemoveAll(filepath.Join("testdata", "logs"))
	_ = os.MkdirAll(filepath.Join("testdata", "logs"), 0755)
	_ = os.WriteFile(filepath.Join("testdata", "logs", ".gitignore"), []byte("*.log"), 0644)
	f, err := NewRotatingFile(RotatingFileOptions{
		Dir:         filepath.Join("testdata", "logs"),
		Filename:    "test",
		MaxFileSize: 10,
	})
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
	f, err = NewRotatingFile(RotatingFileOptions{
		Dir:          filepath.Join("testdata", "logs"),
		Filename:     "test-maxcount",
		MaxFileSize:  10,
		MaxFileCount: 2,
	})
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	_, err = f.Write([]byte("hello, world, hello, world, hello, world"))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
}
