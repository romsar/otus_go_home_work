package config

import (
	"testing"
	"time"

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
				Kafka: KafkaConfig{
					Brokers:     []string{"kafka:9092"},
					GroupID:     "calendar",
					SenderTopic: "calendar-sender-topic",
				},
				Scheduler: SchedulerConfig{
					Interval:        1 * time.Minute,
					EventLifeInDays: 365,
				},
				Sender: SenderConfig{
					Threads: 3,
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

			got, err := Load(tt.args.path)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
