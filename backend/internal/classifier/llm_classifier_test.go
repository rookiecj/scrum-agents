package classifier

import (
	"fmt"
	"testing"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// mockLLMClient is a test double for LLMClient.
type mockLLMClient struct {
	response string
	err      error
}

func (m *mockLLMClient) Complete(prompt string) (string, error) {
	return m.response, m.err
}

func TestLLMClassifier_Classify(t *testing.T) {
	tests := []struct {
		name        string
		llmResponse string
		llmErr      error
		wantPrimary model.ContentCategory
		wantConf    float64
		wantErr     bool
	}{
		{
			name:        "principle classification",
			llmResponse: `{"primary":"원리소개","confidence":0.92,"secondary":"기술소개","secondary_confidence":0.45}`,
			wantPrimary: model.CategoryPrinciple,
			wantConf:    0.92,
		},
		{
			name:        "review classification",
			llmResponse: `{"primary":"사용기","confidence":0.88,"secondary":"생각정리","secondary_confidence":0.3}`,
			wantPrimary: model.CategoryReview,
			wantConf:    0.88,
		},
		{
			name:        "tutorial classification",
			llmResponse: `{"primary":"튜토리얼","confidence":0.95}`,
			wantPrimary: model.CategoryTutorial,
			wantConf:    0.95,
		},
		{
			name:        "news classification",
			llmResponse: `{"primary":"뉴스/분석","confidence":0.78,"secondary":"기술소개","secondary_confidence":0.55}`,
			wantPrimary: model.CategoryNews,
			wantConf:    0.78,
		},
		{
			name:        "opinion classification",
			llmResponse: `{"primary":"생각정리","confidence":0.85}`,
			wantPrimary: model.CategoryOpinion,
			wantConf:    0.85,
		},
		{
			name:        "tech intro classification",
			llmResponse: `{"primary":"기술소개","confidence":0.90}`,
			wantPrimary: model.CategoryTechIntro,
			wantConf:    0.90,
		},
		{
			name:        "invalid category",
			llmResponse: `{"primary":"unknown","confidence":0.9}`,
			wantErr:     true,
		},
		{
			name:        "invalid json",
			llmResponse: `not json`,
			wantErr:     true,
		},
		{
			name:    "llm error",
			llmErr:  fmt.Errorf("API error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockLLMClient{response: tt.llmResponse, err: tt.llmErr}
			classifier := NewLLMClassifier(client)

			result, err := classifier.Classify("some content")

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Primary != tt.wantPrimary {
				t.Errorf("primary = %q, want %q", result.Primary, tt.wantPrimary)
			}
			if result.Confidence != tt.wantConf {
				t.Errorf("confidence = %f, want %f", result.Confidence, tt.wantConf)
			}
		})
	}
}

func TestClassificationPrompt(t *testing.T) {
	content := "This is a test article about how TCP works."
	prompt := ClassificationPrompt(content)

	if prompt == "" {
		t.Error("expected non-empty prompt")
	}

	// Check that all categories are mentioned
	categories := model.AllCategories()
	for _, cat := range categories {
		if !containsString(prompt, string(cat)) {
			t.Errorf("prompt should contain category %q", cat)
		}
	}

	// Check that content is included
	if !containsString(prompt, content) {
		t.Error("prompt should contain the content")
	}
}

func TestClassificationPrompt_Truncation(t *testing.T) {
	// Create content longer than 4000 chars
	longContent := make([]byte, 5000)
	for i := range longContent {
		longContent[i] = 'a'
	}

	prompt := ClassificationPrompt(string(longContent))
	// Content should be truncated + "..." appended
	if !containsString(prompt, "...") {
		t.Error("long content should be truncated with ...")
	}
}

func TestIsValidCategory(t *testing.T) {
	tests := []struct {
		cat  model.ContentCategory
		want bool
	}{
		{model.CategoryPrinciple, true},
		{model.CategoryReview, true},
		{model.CategoryOpinion, true},
		{model.CategoryTechIntro, true},
		{model.CategoryTutorial, true},
		{model.CategoryNews, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.cat), func(t *testing.T) {
			got := isValidCategory(tt.cat)
			if got != tt.want {
				t.Errorf("isValidCategory(%q) = %v, want %v", tt.cat, got, tt.want)
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
