package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_loadConfig(t *testing.T) {
	type args struct {
		path string
		envs map[string]string
	}

	defaultEnvs := map[string]string{
		"LOG_LEVEL": "debug",
		"DB_DRIVER": "inmemory",
	}

	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "log config",
			args: args{
				envs: map[string]string{
					"LOG_LEVEL": "error",
					"DB_DRIVER": "postgres",
				},
			},
			want: &Config{
				Log: LogConfig{
					Level: "error",
				},
				DBDriver: "postgres",
				REST:     RESTConfig{Address: ":8080"},
				GRPC:     GRPCConfig{Address: ":8081"},
				PostgreSQL: PostgreSQLConfig{
					Host:     "",
					Port:     0,
					User:     "",
					Password: "",
					Database: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, val := range defaultEnvs {
				t.Setenv(key, val)
			}
			for key, val := range tt.args.envs {
				t.Setenv(key, val)
			}

			got, err := loadConfig(tt.args.path)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
