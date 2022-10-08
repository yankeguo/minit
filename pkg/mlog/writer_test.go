package mlog

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestLoggerWriter(t *testing.T) {
	out := &bytes.Buffer{}
	l := log.New(out, "aaa", log.Lshortfile)
	w := NewLoggerWriter(l, "bbb ")
	_, err := w.Write([]byte("hello,world\nbbb"))
	require.NoError(t, err)
	err = w.Close()
	require.NoError(t, err)
	require.Equal(t, "aaawriter_test.go:14: bbb hello,world\naaawriter_test.go:16: bbb bbb\n", out.String())
}
