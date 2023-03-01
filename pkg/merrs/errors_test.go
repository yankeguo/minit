package merrs

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewErrorGroup(t *testing.T) {
	eg := NewErrorGroup()
	eg.Add(errors.New("hello"))
	eg.Add(nil)
	eg.Add(errors.New("world"))
	require.Equal(t, "#0: hello; #2: world", eg.Unwrap().Error())

	eg = NewErrorGroup()
	eg.Add(nil)
	eg.Add(nil)
	require.NoError(t, eg.Unwrap())

	eg.Set(3, errors.New("BBB"))
	require.Error(t, eg.Unwrap())

	errs := eg.Unwrap().(Errors)
	require.Equal(t, 4, len(errs))
	require.NoError(t, errs[0])
	require.NoError(t, errs[1])
	require.NoError(t, errs[2])
	require.Error(t, errs[3])
}
