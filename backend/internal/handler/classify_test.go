package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

type mockClassifier struct {
	result *model.ClassificationResult
	err    error
}

func (m *mockClassifier) Classify(content string) (*model.ClassificationResult, error) {
	return m.result, m.err
}

func TestHandleClassify(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		result     *model.ClassificationResult
		classErr   error
		wantStatus int
		wantErr    bool
	}{
		{
			name: "successful classification",
			body: `{"content":"This article explains how TCP works..."}`,
			result: &model.ClassificationResult{
				Primary:    model.CategoryPrinciple,
				Confidence: 0.92,
			},
			wantStatus: 200,
		},
		{
			name:       "empty content",
			body:       `{"content":""}`,
			wantStatus: 400,
			wantErr:    true,
		},
		{
			name:       "missing content",
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
		{
			name:       "classifier error",
			body:       `{"content":"some content"}`,
			classErr:   fmt.Errorf("LLM unavailable"),
			wantStatus: 500,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cls := &mockClassifier{result: tt.result, err: tt.classErr}
			handler := HandleClassify(cls, nil)

			req := httptest.NewRequest("POST", "/api/classify", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var resp ClassifyResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if tt.wantErr {
				if resp.Error == "" {
					t.Error("expected error in response")
				}
				return
			}

			if resp.Classification == nil {
				t.Fatal("expected classification result")
			}
			if resp.Classification.Primary != tt.result.Primary {
				t.Errorf("primary = %q, want %q", resp.Classification.Primary, tt.result.Primary)
			}
		})
	}
}
