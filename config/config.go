package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	PostgresDSN     string
	RedisAddr       string
	SlackWebhookURL string
	APPDomain       string
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() *Config {
	viper.SetConfigName("config")    // Name of the config file (without extension)
	viper.SetConfigType("yaml")      // Format of the config file (yaml, json, etc.)
	viper.AddConfigPath("./config/") // Path to look for the config file in

	// Read from environment variables as well
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("POSTGRES_DSN", "postgres://user:password@localhost:5432/dbname?sslmode=disable")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, using defaults: %v", err)
	}

	// Return the loaded configuration
	return &Config{
		PostgresDSN:     viper.GetString("POSTGRES_DSN"),
		RedisAddr:       viper.GetString("REDIS_ADDR"),
		SlackWebhookURL: viper.GetString("SLACK_WEBHOOK_URL"),
		APPDomain:       viper.GetString("APP_DOMAIN"),
	}
}
