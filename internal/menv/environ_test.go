package menv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnviron(t *testing.T) {
	m := Environ()
	require.NotEmpty(t, m["SHELL"])
}
