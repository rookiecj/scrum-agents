package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGeminiProvider_Complete(t *testing.T) {
	tests := []struct {
		name       string
		response   geminiResponse
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful response",
			response: geminiResponse{
				Candidates: []struct {
					Content struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					} `json:"content"`
				}{
					{
						Content: struct {
							Parts []struct {
								Text string `json:"text"`
							} `json:"parts"`
						}{
							Parts: []struct {
								Text string `json:"text"`
							}{{Text: "Hello from Gemini"}},
						},
					},
				},
			},
			statusCode: 200,
			want:       "Hello from Gemini",
		},
		{
			name: "api error",
			response: geminiResponse{
				Error: &struct {
					Message string `json:"message"`
					Code    int    `json:"code"`
				}{Message: "API key not valid", Code: 400},
			},
			statusCode: 400,
			wantErr:    true,
		},
		{
			name: "empty candidates",
			response: geminiResponse{
				Candidates: []struct {
					Content struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					} `json:"content"`
				}{},
			},
			statusCode: 200,
			wantErr:    true,
		},
		{
			name: "empty parts",
			response: geminiResponse{
				Candidates: []struct {
					Content struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					} `json:"content"`
				}{
					{
						Content: struct {
							Parts []struct {
								Text string `json:"text"`
							} `json:"parts"`
						}{
							Parts: []struct {
								Text string `json:"text"`
							}{},
						},
					},
				},
			},
			statusCode: 200,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify API key header
				if r.Header.Get("x-goog-api-key") != "test-key" {
					t.Error("missing x-goog-api-key header")
				}
				// Verify URL contains model name
				if r.URL.Path != "/models/gemini-2.0-flash:generateContent" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			p := &GeminiProvider{
				config:  Config{APIKey: "test-key", Model: "gemini-2.0-flash", MaxTokens: 4096},
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

func TestGeminiProvider_Name(t *testing.T) {
	p := NewGeminiProvider(Config{APIKey: "test", Timeout: 5 * time.Second})
	if p.Name() != ProviderGemini {
		t.Errorf("Name() = %q, want %q", p.Name(), ProviderGemini)
	}
}

func TestGeminiProvider_ImplementsProvider(t *testing.T) {
	var _ Provider = &GeminiProvider{}
}

func TestGeminiProvider_RequestFormat(t *testing.T) {
	var receivedBody geminiRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)

		resp := geminiResponse{
			Candidates: []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			}{
				{
					Content: struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					}{
						Parts: []struct {
							Text string `json:"text"`
						}{{Text: "response"}},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := &GeminiProvider{
		config:  Config{APIKey: "test-key", Model: "gemini-2.0-flash", MaxTokens: 4096},
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := p.Complete("hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify request body format matches Gemini API spec
	if len(receivedBody.Contents) != 1 {
		t.Fatalf("expected 1 content, got %d", len(receivedBody.Contents))
	}
	if len(receivedBody.Contents[0].Parts) != 1 {
		t.Fatalf("expected 1 part, got %d", len(receivedBody.Contents[0].Parts))
	}
	if receivedBody.Contents[0].Parts[0].Text != "hello world" {
		t.Errorf("prompt = %q, want %q", receivedBody.Contents[0].Parts[0].Text, "hello world")
	}
}
