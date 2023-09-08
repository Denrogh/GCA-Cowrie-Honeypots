package config

import (
	"os"
	"strconv"
)

type Config struct {
	LogPath            string `env:"LOGPATH"`
	EnablePolling      bool   `env:"ENABLE_POLLING"`
	PollingInterval    string `env:"POLLING_INTERVAL"`
	KafkaTopic         string `env:"KAFKA_TOPIC"`
	KafkaBrokers       string `env:"KAFKA_BROKERS"`
	OpensearchIndex    string `env:"INDEX"`
	OpensearchUrl      string `env:"OPENSEARCH_URL"`
	OpensearchUserName string `env:"OPENSEARCH_USER"`
	OpensearchPassword string `env:"OPENSEARCH_PASSWORD"`
}

func LoadConfig() (*Config, error) {

	// Read the environment variables
	config := Config{
		LogPath:            os.Getenv("LOGPATH"),
		EnablePolling:      os.Getenv("ENABLE_POLLING") == "true",
		PollingInterval:    os.Getenv("POLLING_INTERVAL"),
		KafkaTopic:         os.Getenv("KAFKA_TOPIC"),
		KafkaBrokers:       os.Getenv("KAFKA_BROKERS"),
		OpensearchIndex:    os.Getenv("INDEX"),
		OpensearchUrl:      os.Getenv("OPENSEARCH_URL"),
		OpensearchUserName: os.Getenv("OPENSEARCH_USER"),
		OpensearchPassword: os.Getenv("OPENSEARCH_PASSWORD"),
	}

	return &config, nil
}

func parseIntegerEnv(value string) int {
	// Parse the integer value from the environment variable
	// You can add error handling if necessary
	if value == "" {
		return 0
	}

	result, _ := strconv.Atoi(value)
	return result
}
