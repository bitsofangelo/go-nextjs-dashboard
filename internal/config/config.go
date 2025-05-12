package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppPort     string `mapstructure:"APP_PORT"`
	AppTimezone string `mapstructure:"APP_TIMEZONE"`

	DBHost string `mapstructure:"DB_HOST"`
	DBPort int    `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	DBName string `mapstructure:"DB_NAME"`

	LogLevel  string `mapstructure:"LOG_LEVEL"`  // "debug" | "info" | "warn" | "error"
	LogFormat string `mapstructure:"LOG_FORMAT"` // "text" | "json"
	LogPath   string `mapstructure:"LOG_PATH"`
	LogOutput string `mapstructure:"LOG_OUTPUT"` // "stdout" (default) | "./logs/app.log" | "/var/log/<app_name>/app.log

	// Deprecated TODO: to be removed
	DatabaseURL string `mapstructure:"MYSQL_URI"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
}

// Cfg Deprecated
// TODO: to be removed
var Cfg Config

// LoadConfig Deprecated
// TODO: to be removed
func LoadConfig() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(fmt.Errorf("error loading .env file: %w", err))
	}

	Cfg = Config{
		DatabaseURL: os.Getenv("MYSQL_URI"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading .env: %w", err)
	}

	// Also allow real ENV to override
	viper.AutomaticEnv()

	// Unmarshal into your struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
