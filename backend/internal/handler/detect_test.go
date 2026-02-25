package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleDetect(t *testing.T) {
	handler := HandleDetect()

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantType   string
		wantErr    bool
	}{
		{
			name:       "detect article",
			body:       `{"url":"https://example.com/blog/post"}`,
			wantStatus: 200,
			wantType:   "article",
		},
		{
			name:       "detect youtube",
			body:       `{"url":"https://www.youtube.com/watch?v=abc"}`,
			wantStatus: 200,
			wantType:   "youtube",
		},
		{
			name:       "detect pdf",
			body:       `{"url":"https://example.com/paper.pdf"}`,
			wantStatus: 200,
			wantType:   "pdf",
		},
		{
			name:       "missing url",
			body:       `{}`,
			wantStatus: 400,
			wantErr:    true,
		},
		{
			name:       "invalid json",
			body:       `not json`,
			wantStatus: 400,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/detect", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var resp DetectResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if tt.wantErr {
				if resp.Error == "" {
					t.Error("expected error in response")
				}
				return
			}

			if string(resp.LinkInfo.LinkType) != tt.wantType {
				t.Errorf("link_type = %q, want %q", resp.LinkInfo.LinkType, tt.wantType)
			}
		})
	}
}

func TestHandleExtract(t *testing.T) {
	// Create a test server that serves HTML
	htmlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head><title>Test Page</title></head><body><p>Hello from test server.</p></body></html>`))
	}))
	defer htmlServer.Close()

	handler := HandleExtract()

	t.Run("extract article content", func(t *testing.T) {
		body := `{"url":"` + htmlServer.URL + `"}`
		req := httptest.NewRequest("POST", "/api/extract", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != 200 {
			t.Fatalf("status = %d, want 200", rec.Code)
		}

		var resp ExtractResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.LinkInfo.Title != "Test Page" {
			t.Errorf("title = %q, want %q", resp.LinkInfo.Title, "Test Page")
		}
		if resp.Content == "" {
			t.Error("expected non-empty content")
		}
	})

	t.Run("unsupported type returns info", func(t *testing.T) {
		body := `{"url":"https://www.youtube.com/watch?v=abc"}`
		req := httptest.NewRequest("POST", "/api/extract", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != 200 {
			t.Fatalf("status = %d, want 200", rec.Code)
		}

		var resp ExtractResponse
		json.NewDecoder(rec.Body).Decode(&resp)

		if resp.Error == "" {
			t.Error("expected error message for unsupported type")
		}
	})

	t.Run("missing url", func(t *testing.T) {
		body := `{}`
		req := httptest.NewRequest("POST", "/api/extract", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != 400 {
			t.Errorf("status = %d, want 400", rec.Code)
		}
	})
}
