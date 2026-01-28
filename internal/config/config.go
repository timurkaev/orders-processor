package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// Kafka
	KafkaBootstrapServers string // Kafka brokers addresses
	KafkaTopic            string // Topic name
	KafkaGroupID          string // Consumer group  ID

	// PostgreSQL
	PostgresHost     string // DB Host (localhost)
	PostgresPort     int    // DB Port
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	LogLever string // (info, debug, error)
	Port     int
}

func Load() (*Config, error) {
	cfg := &Config{
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		KafkaTopic:            getEnv("KAFKA_TOPIC", "orders"),
		KafkaGroupID:          getEnv("KAFKA_GROUP_ID", "order-processor-group"),

		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:       getEnv("POSTGRES_DB", "orders_db"),

		LogLever: getEnv("LOG_LEVER", "info"),
	}

	var err error

	cfg.PostgresPort, err = getEnvAsInt("POSTGRES_PORT", 5432)
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_PORT: %w", err)
	}

	cfg.Port, err = getEnvAsInt("PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.KafkaBootstrapServers == "" {
		return fmt.Errorf("KAFKA_BOOTSTRAP_SERVERS is required")
	}
	if c.KafkaTopic == "" {
		return fmt.Errorf("KAFKA_TOPIC is required")
	}
	if c.PostgresHost == "" {
		return fmt.Errorf("POSTGRES_DB is required")
	}
	if c.PostgresDB == "" {
		return fmt.Errorf("POSTGRES_DB is required")
	}

	if c.PostgresPort < 1 || c.PostgresPort > 65535 {
		return fmt.Errorf("POSTGRES_PORT must be between 1 and 65535")
	}
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("PORT must be between 1 and 65535")
	}

	return nil
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
	)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvAsInt(key string, defaultValue int) (int, error) {
	valueStr := os.Getenv(key)

	if valueStr == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}

	return value, nil
}
