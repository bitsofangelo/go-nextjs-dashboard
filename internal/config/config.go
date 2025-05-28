package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Env = int

const (
	Production Env = iota
	Staging
	Local
)

type Config struct {
	AppEnv      Env
	AppPort     string `mapstructure:"APP_PORT"`
	AppTimezone string `mapstructure:"APP_TIMEZONE"`
	AppDebug    bool   `mapstructure:"APP_DEBUG"`

	DBHost string `mapstructure:"DB_HOST"`
	DBPort int    `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	DBName string `mapstructure:"DB_NAME"`

	LogLevel  string `mapstructure:"LOG_LEVEL"`  // "debug" | "info" | "warn" | "error"
	LogFormat string `mapstructure:"LOG_FORMAT"` // "text" | "json"
	LogPath   string `mapstructure:"LOG_PATH"`
	LogOutput string `mapstructure:"LOG_OUTPUT"` // "stdout" (default) | "./logs/app.log" | "/var/log/<app_name>/app.log

	JWTHmacKey string `mapstructure:"JWT_HMAC_KEY"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading .env: %w", err)
	}

	if len(viper.AllKeys()) == 0 {
		return nil, fmt.Errorf("no .env variables found")
	}

	// Also allow real ENV to override
	viper.AutomaticEnv()

	// Unmarshal into your struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	switch viper.GetString("APP_ENV") {
	case "production":
		cfg.AppEnv = Production
	case "staging":
		cfg.AppEnv = Staging
	default:
		cfg.AppEnv = Local
	}

	return &cfg, nil
}
