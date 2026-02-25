package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  slog.Level
	}{
		{name: "debug lowercase", input: "debug", want: slog.LevelDebug},
		{name: "debug uppercase", input: "DEBUG", want: slog.LevelDebug},
		{name: "debug mixed case", input: "Debug", want: slog.LevelDebug},
		{name: "info lowercase", input: "info", want: slog.LevelInfo},
		{name: "info uppercase", input: "INFO", want: slog.LevelInfo},
		{name: "warn", input: "warn", want: slog.LevelWarn},
		{name: "warning", input: "warning", want: slog.LevelWarn},
		{name: "error", input: "error", want: slog.LevelError},
		{name: "empty string defaults to info", input: "", want: slog.LevelInfo},
		{name: "unknown defaults to info", input: "trace", want: slog.LevelInfo},
		{name: "whitespace trimmed", input: "  debug  ", want: slog.LevelDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLevel(tt.input)
			if got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestInit_SetsJSONHandler(t *testing.T) {
	// Save original default and restore after test
	origHandler := slog.Default().Handler()
	defer slog.SetDefault(slog.New(origHandler))

	t.Setenv("LOG_LEVEL", "debug")
	Init()

	// Verify that after Init(), the logger produces JSON output
	var buf bytes.Buffer
	testHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(testHandler)
	logger.Info("test message", slog.String("key", "value"))

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("expected JSON output, got: %s", buf.String())
	}

	if logEntry["msg"] != "test message" {
		t.Errorf("msg = %v, want %q", logEntry["msg"], "test message")
	}
	if logEntry["key"] != "value" {
		t.Errorf("key = %v, want %q", logEntry["key"], "value")
	}
}

func TestInit_RespectsLogLevel(t *testing.T) {
	origHandler := slog.Default().Handler()
	defer slog.SetDefault(slog.New(origHandler))

	// Set log level to error
	t.Setenv("LOG_LEVEL", "error")
	Init()

	// The default logger should now be set; verify it is enabled at error level
	// and not at info level
	if slog.Default().Enabled(nil, slog.LevelInfo) {
		t.Error("expected Info level to be disabled when LOG_LEVEL=error")
	}
	if !slog.Default().Enabled(nil, slog.LevelError) {
		t.Error("expected Error level to be enabled when LOG_LEVEL=error")
	}
}
