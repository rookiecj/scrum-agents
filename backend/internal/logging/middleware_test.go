package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware_LogsRequestInfo(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(slog.Default())

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Middleware(inner)

	req := httptest.NewRequest("GET", "/api/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("expected JSON log output, got: %s", buf.String())
	}

	if logEntry["msg"] != "http request" {
		t.Errorf("msg = %v, want %q", logEntry["msg"], "http request")
	}
	if logEntry["method"] != "GET" {
		t.Errorf("method = %v, want %q", logEntry["method"], "GET")
	}
	if logEntry["path"] != "/api/test" {
		t.Errorf("path = %v, want %q", logEntry["path"], "/api/test")
	}
	if status, ok := logEntry["status"].(float64); !ok || int(status) != 200 {
		t.Errorf("status = %v, want 200", logEntry["status"])
	}
	if _, ok := logEntry["duration_ms"]; !ok {
		t.Error("expected duration_ms in log output")
	}
}

func TestMiddleware_CapturesStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{name: "200 OK", statusCode: http.StatusOK},
		{name: "400 Bad Request", statusCode: http.StatusBadRequest},
		{name: "404 Not Found", statusCode: http.StatusNotFound},
		{name: "500 Internal Server Error", statusCode: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
			slog.SetDefault(slog.New(handler))

			inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})

			middleware := Middleware(inner)

			req := httptest.NewRequest("POST", "/api/classify", nil)
			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, req)

			var logEntry map[string]any
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				t.Fatalf("expected JSON log output, got: %s", buf.String())
			}

			if status, ok := logEntry["status"].(float64); !ok || int(status) != tt.statusCode {
				t.Errorf("status = %v, want %d", logEntry["status"], tt.statusCode)
			}
		})
	}
}

func TestMiddleware_DefaultStatusCode(t *testing.T) {
	// If the handler writes directly without calling WriteHeader,
	// the default status should be 200.
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	middleware := Middleware(inner)

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("expected JSON log output, got: %s", buf.String())
	}

	if status, ok := logEntry["status"].(float64); !ok || int(status) != 200 {
		t.Errorf("status = %v, want 200", logEntry["status"])
	}
}
