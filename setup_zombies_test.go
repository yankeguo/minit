//go:build linux
// +build linux

package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheckProcStatIsZombie(t *testing.T) {
	var res bool
	res = checkProcStatIsZombie([]byte("299923 (kworker/2:1-cgroup_pidlist_destroy) R 2 0 0 0 -1 69238880 0 0 0 0 9 153 0 0 20 0 1 0 78232531 0 0 18446744073709551615 0 0 0 0 0 0 0 2147483647 0 0 0 0 17 2 0 0 0 0 0 0 0 0 0 0 0 0 0"))
	require.False(t, res)
	res = checkProcStatIsZombie([]byte("299923 (kworker/2:1-cgroup_pidlist_destroy) Z 2 0 0 0 -1 69238880 0 0 0 0 9 153 0 0 20 0 1 0 78232531 0 0 18446744073709551615 0 0 0 0 0 0 0 2147483647 0 0 0 0 17 2 0 0 0 0 0 0 0 0 0 0 0 0 0"))
	require.True(t, res)
}
