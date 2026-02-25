package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
	"github.com/rookiecj/scrum-agents/backend/internal/summarizer"
)

type mockSummarizerLLM struct {
	response string
	err      error
}

func (m *mockSummarizerLLM) Complete(prompt string) (string, error) {
	return m.response, m.err
}

func newTestSummarizer(t *testing.T) *summarizer.Summarizer {
	t.Helper()
	// Create minimal templates for testing
	reg := newTestRegistry(t)
	return summarizer.NewSummarizer(reg, 0.6)
}

func newTestRegistry(t *testing.T) *summarizer.TemplateRegistry {
	t.Helper()
	// Try to load from the actual prompts directory
	candidates := []string{
		"../../../prompts",
		"../../prompts",
	}
	for _, dir := range candidates {
		reg, err := summarizer.LoadTemplates(dir)
		if err == nil {
			return reg
		}
	}
	t.Skip("prompts directory not found")
	return nil
}

func TestHandleSummarize(t *testing.T) {
	s := newTestSummarizer(t)
	client := &mockSummarizerLLM{response: "## 핵심 원리\nTCP is..."}

	handler := HandleSummarize(s, client)

	tests := []struct {
		name       string
		body       any
		wantStatus int
		wantError  bool
	}{
		{
			name: "valid request with high confidence",
			body: SummarizeRequest{
				Content: "TCP works by establishing connections...",
				Classification: &model.ClassificationResult{
					Primary:    model.CategoryPrinciple,
					Confidence: 0.92,
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "valid request with low confidence",
			body: SummarizeRequest{
				Content: "Some ambiguous content...",
				Classification: &model.ClassificationResult{
					Primary:    model.CategoryPrinciple,
					Confidence: 0.3,
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "empty content",
			body:       SummarizeRequest{Content: "", Classification: &model.ClassificationResult{Primary: model.CategoryPrinciple, Confidence: 0.9}},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "missing classification",
			body:       SummarizeRequest{Content: "some content"},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "invalid body",
			body:       "not json",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/summarize", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			var resp SummarizeResponse
			json.NewDecoder(w.Body).Decode(&resp)

			if tt.wantError && resp.Error == "" {
				t.Error("expected error in response")
			}
			if !tt.wantError && resp.Error != "" {
				t.Errorf("unexpected error: %s", resp.Error)
			}
			if !tt.wantError && resp.Result == nil {
				t.Error("expected result in response")
			}
		})
	}
}

func TestHandleSummarize_LLMError(t *testing.T) {
	s := newTestSummarizer(t)
	client := &mockSummarizerLLM{err: fmt.Errorf("API down")}

	handler := HandleSummarize(s, client)

	body := SummarizeRequest{
		Content: "test content",
		Classification: &model.ClassificationResult{
			Primary:    model.CategoryPrinciple,
			Confidence: 0.9,
		},
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/summarize", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
