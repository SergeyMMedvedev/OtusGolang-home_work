package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	testDir := "./testdata/env"
	expected := Environment{
		"BAR": EnvValue{
			"bar",
			false,
		},
		"EMPTY": EnvValue{
			"",
			false,
		},
		"FOO": EnvValue{
			"   foo\nwith new line",
			false,
		},
		"HELLO": EnvValue{
			"\"hello\"",
			false,
		},
	}
	env, err := ReadDir(testDir)
	require.NoError(t, err)
	require.Equal(t, expected, env)
	_, ok := os.LookupEnv("UNSET")
	require.False(t, ok)
}

func TestReadDirReplace(t *testing.T) {
	testDir := "./testdata/env"
	expected := Environment{
		"BAR": EnvValue{
			"bar",
			true,
		},
		"EMPTY": EnvValue{
			"",
			false,
		},
		"FOO": EnvValue{
			"   foo\nwith new line",
			false,
		},
		"HELLO": EnvValue{
			"\"hello\"",
			true,
		},
	}
	os.Setenv("BAR", "bar")
	os.Setenv("HELLO", "hello")
	env, err := ReadDir(testDir)
	require.NoError(t, err)
	require.Equal(t, expected, env)
}

func TestReadDirIgnoreEquals(t *testing.T) {
	testDir := "./testdata/env2"
	expected := Environment{
		"FOO": EnvValue{
			"   foo\nwith new line",
			false,
		},
	}
	env, err := ReadDir(testDir)
	require.NoError(t, err)
	require.Equal(t, expected, env)
}
