package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	type args struct {
		cmd  []string
		envs Environment
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"basic test", args{
			cmd:  []string{"true"},
			envs: Environment{},
		}, 0},
		{"bad code test", args{
			cmd:  []string{"-notvalid"},
			envs: Environment{},
		}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, RunCmd(tt.args.cmd, tt.args.envs))
		})
	}
}
