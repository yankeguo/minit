package mtmpl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	buf, err := Execute(`{{$x := 1}}{{add $x 1}}{{"\n"}}{{stringsToUpper .A}}`, map[string]interface{}{"A": "B"})
	require.NoError(t, err)
	require.Equal(t, "2\nB", strings.TrimSpace(string(buf)))
}
