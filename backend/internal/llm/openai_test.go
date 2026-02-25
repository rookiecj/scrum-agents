package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOpenAIProvider_Complete(t *testing.T) {
	tests := []struct {
		name       string
		response   openaiResponse
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful response",
			response: openaiResponse{
				Choices: []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				}{{Message: struct {
					Content string `json:"content"`
				}{Content: "Hello from OpenAI"}}},
			},
			statusCode: 200,
			want:       "Hello from OpenAI",
		},
		{
			name: "api error",
			response: openaiResponse{
				Error: &struct {
					Message string `json:"message"`
				}{Message: "rate limit exceeded"},
			},
			statusCode: 429,
			wantErr:    true,
		},
		{
			name: "empty choices",
			response: openaiResponse{
				Choices: []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				}{},
			},
			statusCode: 200,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify auth header
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Error("missing Authorization header")
				}

				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			p := &OpenAIProvider{
				config:  Config{APIKey: "test-key", Model: "gpt-4o", MaxTokens: 1024},
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

func TestOpenAIProvider_Name(t *testing.T) {
	p := NewOpenAIProvider(Config{APIKey: "test", Timeout: 5 * time.Second})
	if p.Name() != ProviderOpenAI {
		t.Errorf("Name() = %q, want %q", p.Name(), ProviderOpenAI)
	}
}

func TestOpenAIProvider_ImplementsProvider(t *testing.T) {
	var _ Provider = &OpenAIProvider{}
}
