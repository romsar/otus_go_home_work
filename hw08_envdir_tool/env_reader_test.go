package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		want Environment
	}{
		{"basic test", "./testdata/env", Environment{
			"BAR": EnvValue{
				Value:      "bar",
				NeedRemove: false,
			},
			"EMPTY": EnvValue{
				Value:      "",
				NeedRemove: true,
			},
			"FOO": EnvValue{
				Value:      "   foo\nwith new line",
				NeedRemove: false,
			},
			"HELLO": EnvValue{
				Value:      "\"hello\"",
				NeedRemove: false,
			},
			"UNSET": EnvValue{
				Value:      "",
				NeedRemove: true,
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			environment, err := ReadDir(tt.dir)
			require.NoError(t, err)

			require.Equal(t, environment, tt.want)
		})
	}
}

func TestReadDirNotFound(t *testing.T) {
	environment, err := ReadDir("./testdata/foobar")
	require.Error(t, err)
	require.Nil(t, environment)
}
