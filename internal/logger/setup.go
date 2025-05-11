package logger

import (
	"log/slog"
	"os"
	"strings"

	"go-nextjs-dashboard/internal/config"
)

func New(cfg *config.Config) *slog.Logger {
	var h slog.Handler

	level := slog.LevelInfo
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	if cfg.LogFormat == "json" {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	return slog.New(h)
}
