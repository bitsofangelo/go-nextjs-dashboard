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

	HashingDriver string `mapstructure:"HASHING_DRIVER"`

	LogLevel  string `mapstructure:"LOG_LEVEL"`  // "debug" | "info" | "warn" | "error"
	LogFormat string `mapstructure:"LOG_FORMAT"` // "text" | "json"
	LogPath   string `mapstructure:"LOG_PATH"`   // "./logs/app.log" | "/var/log/<app_name>/app.log"
	LogOutput string `mapstructure:"LOG_OUTPUT"` // "stdout" (default) | "file"

	JWTHmacKey string `mapstructure:"JWT_HMAC_KEY"`

	MailHost          string `mapstructure:"MAIL_HOST"`
	MailPort          int    `mapstructure:"MAIL_PORT"`
	MailUser          string `mapstructure:"MAIL_USER"`
	MailPass          string `mapstructure:"MAIL_PASS"`
	MailSkipTLSVerify bool   `mapstructure:"MAIL_SKIP_TLS_VERIFY"`
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

	// allow real ENV to override
	viper.AutomaticEnv()

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
