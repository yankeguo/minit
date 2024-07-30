package mtmpl

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	buf, err := Execute(TEST_TMPL, map[string]interface{}{"A": "B"})
	require.NoError(t, err)
	require.Equal(t, "2\nB", strings.TrimSpace(string(buf)))
}
