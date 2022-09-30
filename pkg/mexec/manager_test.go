package mexec

import (
	"os"
	"testing"
)

func TestNewManager(t *testing.T) {
	os.Setenv("AAA-", "BBB")
	t.Log(os.Getenv("AAA-"))
}
