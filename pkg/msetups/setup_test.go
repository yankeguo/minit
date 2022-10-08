package msetups

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetup(t *testing.T) {
	require.Equal(t, 10, setups[0].A)
}
