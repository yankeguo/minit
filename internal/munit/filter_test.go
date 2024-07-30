package munit

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFilterMap(t *testing.T) {
	fm := NewFilterMap("")
	require.NotNil(t, fm)
	require.True(t, fm.Blank())

	fm = NewFilterMap(",,  ,")
	require.NotNil(t, fm)
	require.True(t, fm.Blank())

	fm = NewFilterMap("unit-a,&daemon")
	require.True(t, fm.Match(Unit{
		Name: "unit-a",
	}))
	require.True(t, fm.Match(Unit{
		Name: "unit-b",
		Kind: "daemon",
	}))

	fm = NewFilterMap("unit-a, ,, @group-b, unit-c,,")
	require.NotNil(t, fm)
	require.True(t, fm.Match(Unit{
		Name: "unit-a",
	}))
	require.True(t, fm.Match(Unit{
		Name:  "unit-b",
		Group: "group-b",
	}))
	require.True(t, fm.Match(Unit{
		Name:  "unit-c",
		Group: "group-c",
	}))
	require.False(t, fm.Match(Unit{
		Name:  "unit-d",
		Group: "group-d",
	}))
}

func TestNewFilter(t *testing.T) {
	f := NewFilter("  ,  , , ", ",, ,")
	for i := 0; i < 10; i++ {
		buf := make([]byte, 10)
		rand.Read(buf)
		require.True(t, f.Match(Unit{
			Name:  hex.EncodeToString(buf),
			Group: hex.EncodeToString(buf),
		}))
	}

	f = NewFilter("unit-a,&daemon", "")
	require.False(t, f.Match(Unit{
		Name: "bla",
		Kind: KindCron,
	}))
	require.True(t, f.Match(Unit{
		Name: "bla",
		Kind: KindDaemon,
	}))
	require.True(t, f.Match(Unit{
		Name: "unit-a",
		Kind: KindCron,
	}))

	f = NewFilter("", "unit-a,&daemon")
	require.True(t, f.Match(Unit{
		Name: "bla",
		Kind: KindCron,
	}))
	require.False(t, f.Match(Unit{
		Name: "bla",
		Kind: KindDaemon,
	}))
	require.False(t, f.Match(Unit{
		Name: "unit-a",
		Kind: KindCron,
	}))

	f = NewFilter("", "unit-a,,,@group-c,,")
	require.True(t, f.Match(Unit{
		Name:  "unit-b",
		Group: "group-b",
	}))
	require.False(t, f.Match(Unit{
		Name:  "unit-c",
		Group: "group-c",
	}))

	f = NewFilter("unit-a,,,@group-c,,", "")
	require.False(t, f.Match(Unit{
		Name:  "unit-b",
		Group: "group-b",
	}))
	require.True(t, f.Match(Unit{
		Name:  "unit-c",
		Group: "group-c",
	}))

	f = NewFilter("unit-a,,,@group-c,,", "unit-c2")
	require.False(t, f.Match(Unit{
		Name:  "unit-b",
		Group: "group-b",
	}))
	require.True(t, f.Match(Unit{
		Name:  "unit-c",
		Group: "group-c",
	}))
	require.False(t, f.Match(Unit{
		Name:  "unit-c2",
		Group: "group-c",
	}))
}
