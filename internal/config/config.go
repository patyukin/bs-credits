package config

import (
	"fmt"
	configLoader "github.com/patyukin/mbs-pkg/pkg/config"
)

type Config struct {
	MinLogLevel string `yaml:"min_log_level" validate:"required,oneof=debug info warn error"`
	HTTPServer  struct {
		Port int `yaml:"port" validate:"required,numeric"`
	} `yaml:"http_server" validate:"required"`
	GRPCServer struct {
		Port int `yaml:"port" validate:"required,numeric"`
	} `yaml:"grpc_server" validate:"required"`
	PostgreSQLDSN   string `yaml:"postgresql_dsn" validate:"required"`
	RedisDSN        string `yaml:"redis_dsn" validate:"required"`
	RabbitMQUrl     string `yaml:"rabbitmq_url" validate:"required"`
	TelegramBotName string `yaml:"telegram_bot_name" validate:"required"`
	TracerHost      string `yaml:"tracer_host" validate:"required"`
	Kafka           struct {
		Brokers       []string `yaml:"brokers" validate:"required"`
		Topics        []string `yaml:"topics" validate:"required"`
		ConsumerGroup string   `yaml:"consumer_group" validate:"required"`
	} `yaml:"kafka" validate:"required"`
	GRPC struct {
		AuthService string `yaml:"auth_service" validate:"required"`
	}
}

func LoadConfig() (*Config, error) {
	var config Config
	err := configLoader.LoadConfig(&config)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return &config, nil
}
