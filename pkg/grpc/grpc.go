package grpc

import "time"

type Config struct {
	Port              int16         `yaml:"port" env:"GRPC_SERVER_PORT" env-required:"true"`
	Timeout           time.Duration `yaml:"timeout" env:"GRPC_SERVER_TIMEOUT" env-required:"true"`
	MaxReceiveMsgSize int           `yaml:"maxReceiveMsgSize" env:"GRPC_MAX_RECEIVE_MSG_SIZE" env-default:"12"`
}
