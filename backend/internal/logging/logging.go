// Package logging provides structured logging initialization using log/slog.
//
// It supports JSON-formatted output with configurable log levels via the
// LOG_LEVEL environment variable (debug, info, warn, error). The default
// level is Info.
package logging

import (
	"log/slog"
	"os"
	"strings"
)

// Init initializes the default slog logger with a JSON handler.
// The log level is determined by the LOG_LEVEL environment variable.
// Supported values: "debug", "info", "warn", "error" (case-insensitive).
// If LOG_LEVEL is empty or unrecognized, defaults to Info.
func Init() {
	level := ParseLevel(os.Getenv("LOG_LEVEL"))
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(handler))
}

// ParseLevel converts a string log level name to a slog.Level.
// Supported values: "debug", "info", "warn", "error" (case-insensitive).
// Returns slog.LevelInfo for empty or unrecognized values.
func ParseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
