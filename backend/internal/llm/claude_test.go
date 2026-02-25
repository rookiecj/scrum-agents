package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClaudeProvider_Complete(t *testing.T) {
	tests := []struct {
		name       string
		response   claudeResponse
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful response",
			response: claudeResponse{
				Content: []struct {
					Text string `json:"text"`
				}{{Text: "Hello from Claude"}},
			},
			statusCode: 200,
			want:       "Hello from Claude",
		},
		{
			name: "api error",
			response: claudeResponse{
				Error: &struct {
					Message string `json:"message"`
				}{Message: "invalid API key"},
			},
			statusCode: 401,
			wantErr:    true,
		},
		{
			name: "empty response",
			response: claudeResponse{
				Content: []struct {
					Text string `json:"text"`
				}{},
			},
			statusCode: 200,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify headers
				if r.Header.Get("x-api-key") != "test-key" {
					t.Error("missing x-api-key header")
				}
				if r.Header.Get("anthropic-version") != "2023-06-01" {
					t.Error("missing anthropic-version header")
				}

				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			p := &ClaudeProvider{
				config:  Config{APIKey: "test-key", Model: "claude-sonnet-4-6", MaxTokens: 1024},
				client:  server.Client(),
				baseURL: server.URL,
			}

			got, err := p.Complete("test prompt")
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Complete() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClaudeProvider_Name(t *testing.T) {
	p := NewClaudeProvider(Config{APIKey: "test", Timeout: 5 * time.Second})
	if p.Name() != ProviderClaude {
		t.Errorf("Name() = %q, want %q", p.Name(), ProviderClaude)
	}
}

func TestClaudeProvider_ImplementsProvider(t *testing.T) {
	var _ Provider = &ClaudeProvider{}
}
