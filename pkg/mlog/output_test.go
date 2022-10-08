package mlog

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewWriterOutput(t *testing.T) {
	buf := &bytes.Buffer{}

	o := NewWriterOutput(buf, []byte("a"), []byte("b"))

	_, err := o.Write([]byte("hello\n"))
	require.NoError(t, err)

	_, err = o.ReadFrom(bytes.NewReader([]byte("hello\nworld")))
	require.NoError(t, err)

	require.Equal(t, "ahello\nbahello\nbaworld\nb", buf.String())
}

func TestMultiOutput(t *testing.T) {
	buf1 := &bytes.Buffer{}
	o1 := NewWriterOutput(buf1, []byte("a"), []byte("b"))
	buf2 := &bytes.Buffer{}
	o2 := NewWriterOutput(buf2, []byte("c"), []byte("d"))

	o := MultiOutput(o1, o2)

	_, err := o.Write([]byte("hello\n"))
	require.NoError(t, err)

	_, err = o.ReadFrom(bytes.NewReader([]byte("hello\nworld")))
	require.NoError(t, err)

	require.Equal(t, "ahello\nbahello\nbaworld\nb", buf1.String())
	require.Equal(t, "chello\ndchello\ndcworld\nd", buf2.String())
}
