package summarizer

import (
	"fmt"
	"testing"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

type mockLLMClient struct {
	response string
	err      error
	lastPrompt string
}

func (m *mockLLMClient) Complete(prompt string) (string, error) {
	m.lastPrompt = prompt
	return m.response, m.err
}

func TestSummarizer_Summarize(t *testing.T) {
	dir := findPromptsDir(t)
	reg, err := LoadTemplates(dir)
	if err != nil {
		t.Fatalf("LoadTemplates() error: %v", err)
	}

	tests := []struct {
		name           string
		classification *model.ClassificationResult
		llmResponse    string
		llmErr         error
		wantStyle      string
		wantTemplate   string
		wantLowConf    bool
		wantErr        bool
	}{
		{
			name: "high confidence principle",
			classification: &model.ClassificationResult{
				Primary:    model.CategoryPrinciple,
				Confidence: 0.92,
			},
			llmResponse:  "## 핵심 원리\nTCP works...",
			wantStyle:    "구조적 요약",
			wantTemplate: "원리소개",
		},
		{
			name: "high confidence review",
			classification: &model.ClassificationResult{
				Primary:    model.CategoryReview,
				Confidence: 0.88,
			},
			llmResponse:  "## 장점\nFast...",
			wantStyle:    "평가 중심 요약",
			wantTemplate: "사용기",
		},
		{
			name: "high confidence opinion",
			classification: &model.ClassificationResult{
				Primary:    model.CategoryOpinion,
				Confidence: 0.85,
			},
			llmResponse:  "## 주장\nAI will...",
			wantStyle:    "논점 중심 요약",
			wantTemplate: "생각정리",
		},
		{
			name: "high confidence techintro",
			classification: &model.ClassificationResult{
				Primary:    model.CategoryTechIntro,
				Confidence: 0.90,
			},
			llmResponse:  "## 핵심 기능\nBun 1.0...",
			wantStyle:    "스펙 중심 요약",
			wantTemplate: "기술소개",
		},
		{
			name: "high confidence tutorial",
			classification: &model.ClassificationResult{
				Primary:    model.CategoryTutorial,
				Confidence: 0.95,
			},
			llmResponse:  "## 목표\nLearn React...",
			wantStyle:    "단계별 요약",
			wantTemplate: "튜토리얼",
		},
		{
			name: "high confidence news",
			classification: &model.ClassificationResult{
				Primary:    model.CategoryNews,
				Confidence: 0.78,
			},
			llmResponse:  "## 핵심 사실\nOpenAI...",
			wantStyle:    "팩트 중심 요약",
			wantTemplate: "뉴스/분석",
		},
		{
			name: "low confidence falls back to generic",
			classification: &model.ClassificationResult{
				Primary:    model.CategoryPrinciple,
				Confidence: 0.3,
			},
			llmResponse:  "## 요약\nGeneral summary...",
			wantStyle:    "일반 요약",
			wantTemplate: "generic",
			wantLowConf:  true,
		},
		{
			name: "confidence at threshold uses category template",
			classification: &model.ClassificationResult{
				Primary:    model.CategoryTutorial,
				Confidence: 0.6,
			},
			llmResponse:  "## 목표\nLearn...",
			wantStyle:    "단계별 요약",
			wantTemplate: "튜토리얼",
		},
		{
			name: "confidence just below threshold uses generic",
			classification: &model.ClassificationResult{
				Primary:    model.CategoryTutorial,
				Confidence: 0.59,
			},
			llmResponse:  "## 요약\nGeneral...",
			wantStyle:    "일반 요약",
			wantTemplate: "generic",
			wantLowConf:  true,
		},
		{
			name: "LLM error propagated",
			classification: &model.ClassificationResult{
				Primary:    model.CategoryPrinciple,
				Confidence: 0.92,
			},
			llmErr:  fmt.Errorf("API down"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockLLMClient{response: tt.llmResponse, err: tt.llmErr}
			s := NewSummarizer(reg, 0.6)

			result, err := s.Summarize(client, "test content", tt.classification)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Summary != tt.llmResponse {
				t.Errorf("Summary = %q, want %q", result.Summary, tt.llmResponse)
			}
			if result.Style != tt.wantStyle {
				t.Errorf("Style = %q, want %q", result.Style, tt.wantStyle)
			}
			if result.TemplateUsed != tt.wantTemplate {
				t.Errorf("TemplateUsed = %q, want %q", result.TemplateUsed, tt.wantTemplate)
			}
			if result.LowConfidence != tt.wantLowConf {
				t.Errorf("LowConfidence = %v, want %v", result.LowConfidence, tt.wantLowConf)
			}
			if result.Category != tt.classification.Primary {
				t.Errorf("Category = %q, want %q", result.Category, tt.classification.Primary)
			}
		})
	}
}

func TestSummarizer_SummarizeWithCategory(t *testing.T) {
	dir := findPromptsDir(t)
	reg, err := LoadTemplates(dir)
	if err != nil {
		t.Fatalf("LoadTemplates() error: %v", err)
	}

	tests := []struct {
		name      string
		category  model.ContentCategory
		wantStyle string
	}{
		{"principle", model.CategoryPrinciple, "구조적 요약"},
		{"review", model.CategoryReview, "평가 중심 요약"},
		{"opinion", model.CategoryOpinion, "논점 중심 요약"},
		{"techintro", model.CategoryTechIntro, "스펙 중심 요약"},
		{"tutorial", model.CategoryTutorial, "단계별 요약"},
		{"news", model.CategoryNews, "팩트 중심 요약"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockLLMClient{response: "summary text"}
			s := NewSummarizer(reg, 0.6)

			result, err := s.SummarizeWithCategory(client, "test content", tt.category)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Style != tt.wantStyle {
				t.Errorf("Style = %q, want %q", result.Style, tt.wantStyle)
			}
			if result.Category != tt.category {
				t.Errorf("Category = %q, want %q", result.Category, tt.category)
			}
			if result.LowConfidence {
				t.Error("LowConfidence should be false for direct category summarization")
			}
		})
	}
}

func TestSummarizer_SummarizeWithCategory_LLMError(t *testing.T) {
	dir := findPromptsDir(t)
	reg, err := LoadTemplates(dir)
	if err != nil {
		t.Fatalf("LoadTemplates() error: %v", err)
	}

	client := &mockLLMClient{err: fmt.Errorf("API error")}
	s := NewSummarizer(reg, 0.6)

	_, err = s.SummarizeWithCategory(client, "test", model.CategoryPrinciple)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSummarizer_PromptContainsTemplateInstruction(t *testing.T) {
	dir := findPromptsDir(t)
	reg, err := LoadTemplates(dir)
	if err != nil {
		t.Fatalf("LoadTemplates() error: %v", err)
	}

	client := &mockLLMClient{response: "summary"}
	s := NewSummarizer(reg, 0.6)

	classification := &model.ClassificationResult{
		Primary:    model.CategoryPrinciple,
		Confidence: 0.92,
	}

	_, err = s.Summarize(client, "TCP works by...", classification)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the prompt sent to LLM contains the template instruction keywords
	prompt := client.lastPrompt
	checks := []string{"핵심 원리", "작동 방식", "전제조건", "한계점", "TCP works by"}
	for _, check := range checks {
		if len(prompt) == 0 || !containsStr(prompt, check) {
			t.Errorf("prompt should contain %q", check)
		}
	}
}

func TestNewSummarizer_DefaultThreshold(t *testing.T) {
	dir := findPromptsDir(t)
	reg, _ := LoadTemplates(dir)

	s := NewSummarizer(reg, 0)
	// Should default to 0.6
	client := &mockLLMClient{response: "summary"}

	result, err := s.Summarize(client, "content", &model.ClassificationResult{
		Primary:    model.CategoryPrinciple,
		Confidence: 0.59,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.LowConfidence {
		t.Error("expected low confidence for 0.59 with default 0.6 threshold")
	}
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
