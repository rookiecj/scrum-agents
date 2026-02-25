package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandleProviders(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		wantLen  int
		wantAvail map[string]bool
	}{
		{
			name:    "no keys set",
			envVars: map[string]string{},
			wantLen: 3,
			wantAvail: map[string]bool{
				"claude": false,
				"openai": false,
				"gemini": false,
			},
		},
		{
			name: "claude key set",
			envVars: map[string]string{
				"ANTHROPIC_API_KEY": "test-key",
			},
			wantLen: 3,
			wantAvail: map[string]bool{
				"claude": true,
				"openai": false,
				"gemini": false,
			},
		},
		{
			name: "all keys set",
			envVars: map[string]string{
				"ANTHROPIC_API_KEY": "key1",
				"OPENAI_API_KEY":    "key2",
				"GOOGLE_API_KEY":    "key3",
			},
			wantLen: 3,
			wantAvail: map[string]bool{
				"claude": true,
				"openai": true,
				"gemini": true,
			},
		},
		{
			name: "gemini only",
			envVars: map[string]string{
				"GOOGLE_API_KEY": "gemini-key",
			},
			wantLen: 3,
			wantAvail: map[string]bool{
				"claude": false,
				"openai": false,
				"gemini": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all relevant env vars
			os.Unsetenv("ANTHROPIC_API_KEY")
			os.Unsetenv("OPENAI_API_KEY")
			os.Unsetenv("GOOGLE_API_KEY")

			// Set env vars for this test
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			handler := HandleProviders()
			req := httptest.NewRequest("GET", "/api/providers", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}

			if rec.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Content-Type = %q, want %q", rec.Header().Get("Content-Type"), "application/json")
			}

			var providers []ProviderInfo
			if err := json.NewDecoder(rec.Body).Decode(&providers); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if len(providers) != tt.wantLen {
				t.Fatalf("got %d providers, want %d", len(providers), tt.wantLen)
			}

			for _, p := range providers {
				wantAvail, ok := tt.wantAvail[p.Name]
				if !ok {
					t.Errorf("unexpected provider: %s", p.Name)
					continue
				}
				if p.Available != wantAvail {
					t.Errorf("provider %s: available = %v, want %v", p.Name, p.Available, wantAvail)
				}
			}

			// Cleanup
			for k := range tt.envVars {
				os.Unsetenv(k)
			}
		})
	}
}

func TestHandleProviders_EnvVarNames(t *testing.T) {
	// Verify that each provider has the correct envVar field
	handler := HandleProviders()
	req := httptest.NewRequest("GET", "/api/providers", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var providers []ProviderInfo
	if err := json.NewDecoder(rec.Body).Decode(&providers); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedEnvVars := map[string]string{
		"claude": "ANTHROPIC_API_KEY",
		"openai": "OPENAI_API_KEY",
		"gemini": "GOOGLE_API_KEY",
	}

	for _, p := range providers {
		expected, ok := expectedEnvVars[p.Name]
		if !ok {
			t.Errorf("unexpected provider: %s", p.Name)
			continue
		}
		if p.EnvVar != expected {
			t.Errorf("provider %s: envVar = %q, want %q", p.Name, p.EnvVar, expected)
		}
	}
}

func TestHandleProviders_DoesNotExposeKeys(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "secret-key-value")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	handler := HandleProviders()
	req := httptest.NewRequest("GET", "/api/providers", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	if containsString(body, "secret-key-value") {
		t.Error("response body should not contain the actual API key value")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
