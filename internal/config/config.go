package config

import (
	"fmt"
	configLoader "github.com/patyukin/mbs-pkg/pkg/config"
)

type Config struct {
	MinLogLevel string `yaml:"min_log_level" validate:"required,oneof=debug info warn error"`
	JwtSecret   string `yaml:"jwt_secret" validate:"required"`
	HttpServer  struct {
		Port int `yaml:"port" validate:"required,numeric"`
	} `yaml:"http_server" validate:"required"`
	GRPCServer struct {
		Port              int `yaml:"port" validate:"required,numeric"`
		MaxConnectionIdle int `yaml:"max_connection_idle"`
		Timeout           int `yaml:"timeout"`
		MaxConnectionAge  int `yaml:"max_connection_age"`
	} `yaml:"grpc_server" validate:"required"`
	PostgreSQLDSN   string `yaml:"postgresql_dsn" validate:"required"`
	RedisDSN        string `yaml:"redis_dsn" validate:"required"`
	RabbitMQUrl     string `yaml:"rabbitmq_url" validate:"required"`
	TelegramBotName string `yaml:"telegram_bot_name" validate:"required"`
	TracerHost      string `yaml:"tracer_host" validate:"required"`
	Kafka           struct {
		Brokers []string `yaml:"brokers" validate:"required"`
		Topic   []string `yaml:"topic" validate:"required"`
	} `yaml:"kafka" validate:"required"`
}

func LoadConfig() (*Config, error) {
	var config Config
	err := configLoader.LoadConfig(&config)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return &config, nil
}
