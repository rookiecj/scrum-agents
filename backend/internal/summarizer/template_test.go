package summarizer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

func TestLoadTemplates(t *testing.T) {
	dir := findPromptsDir(t)

	reg, err := LoadTemplates(dir)
	if err != nil {
		t.Fatalf("LoadTemplates() error: %v", err)
	}

	// Verify all 6 categories are loaded
	for _, cat := range model.AllCategories() {
		tmpl := reg.Get(cat)
		if tmpl == nil {
			t.Errorf("missing template for category: %s", cat)
		}
	}

	// Verify generic template exists
	if reg.GetGeneric() == nil {
		t.Error("missing generic template")
	}
}

func TestLoadTemplates_MissingFile(t *testing.T) {
	dir := t.TempDir()

	_, err := LoadTemplates(dir)
	if err == nil {
		t.Error("expected error for missing template files")
	}
}

func TestTemplateRegistry_Validate(t *testing.T) {
	dir := findPromptsDir(t)

	reg, err := LoadTemplates(dir)
	if err != nil {
		t.Fatalf("LoadTemplates() error: %v", err)
	}

	if err := reg.Validate(); err != nil {
		t.Errorf("Validate() error: %v", err)
	}
}

func TestTemplateRegistry_Validate_MissingTemplate(t *testing.T) {
	reg := &TemplateRegistry{
		templates: make(map[model.ContentCategory]*PromptTemplate),
		generic:   &PromptTemplate{Category: "generic", Instruction: "test", Sections: []string{"a"}},
	}

	err := reg.Validate()
	if err == nil {
		t.Error("expected validation error for missing templates")
	}
}

func TestTemplateRegistry_Validate_EmptyInstruction(t *testing.T) {
	reg := &TemplateRegistry{
		templates: make(map[model.ContentCategory]*PromptTemplate),
		generic:   &PromptTemplate{Category: "generic", Instruction: "test", Sections: []string{"a"}},
	}

	for _, cat := range model.AllCategories() {
		reg.templates[cat] = &PromptTemplate{
			Category:    string(cat),
			Instruction: "",
			Sections:    []string{"a"},
		}
	}

	err := reg.Validate()
	if err == nil {
		t.Error("expected validation error for empty instruction")
	}
}

func TestTemplateRegistry_Validate_NoSections(t *testing.T) {
	reg := &TemplateRegistry{
		templates: make(map[model.ContentCategory]*PromptTemplate),
		generic:   &PromptTemplate{Category: "generic", Instruction: "test", Sections: []string{"a"}},
	}

	for _, cat := range model.AllCategories() {
		reg.templates[cat] = &PromptTemplate{
			Category:    string(cat),
			Instruction: "test instruction",
			Sections:    []string{},
		}
	}

	err := reg.Validate()
	if err == nil {
		t.Error("expected validation error for empty sections")
	}
}

func TestTemplateRegistry_Validate_MissingGeneric(t *testing.T) {
	reg := &TemplateRegistry{
		templates: make(map[model.ContentCategory]*PromptTemplate),
	}

	for _, cat := range model.AllCategories() {
		reg.templates[cat] = &PromptTemplate{
			Category:    string(cat),
			Instruction: "test",
			Sections:    []string{"a"},
		}
	}

	err := reg.Validate()
	if err == nil {
		t.Error("expected validation error for missing generic template")
	}
}

func TestTemplateRegistry_Get(t *testing.T) {
	dir := findPromptsDir(t)
	reg, err := LoadTemplates(dir)
	if err != nil {
		t.Fatalf("LoadTemplates() error: %v", err)
	}

	tests := []struct {
		name     string
		category model.ContentCategory
		wantCat  string
	}{
		{"principle", model.CategoryPrinciple, "원리소개"},
		{"review", model.CategoryReview, "사용기"},
		{"opinion", model.CategoryOpinion, "생각정리"},
		{"techintro", model.CategoryTechIntro, "기술소개"},
		{"tutorial", model.CategoryTutorial, "튜토리얼"},
		{"news", model.CategoryNews, "뉴스/분석"},
		{"unknown falls back to generic", model.ContentCategory("unknown"), "generic"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := reg.Get(tt.category)
			if tmpl == nil {
				t.Fatal("expected non-nil template")
			}
			if tmpl.Category != tt.wantCat {
				t.Errorf("Category = %q, want %q", tmpl.Category, tt.wantCat)
			}
		})
	}
}

func TestTemplateRegistry_Categories(t *testing.T) {
	dir := findPromptsDir(t)
	reg, err := LoadTemplates(dir)
	if err != nil {
		t.Fatalf("LoadTemplates() error: %v", err)
	}

	cats := reg.Categories()
	if len(cats) != 6 {
		t.Errorf("got %d categories, want 6", len(cats))
	}
}

func TestPromptTemplate_BuildPrompt(t *testing.T) {
	tmpl := &PromptTemplate{
		Category:    "원리소개",
		Style:       "구조적 요약",
		Sections:    []string{"핵심 원리", "작동 방식"},
		Instruction: "구조적으로 요약하세요.",
	}

	content := "This is a test article about how TCP works."
	prompt := tmpl.BuildPrompt(content)

	// Check prompt contains key elements
	checks := []string{
		"전문 콘텐츠 요약기",
		"구조적 요약",
		"구조적으로 요약하세요.",
		"This is a test article",
		"한국어로 요약하세요",
	}

	for _, check := range checks {
		if !strings.Contains(prompt, check) {
			t.Errorf("prompt should contain %q", check)
		}
	}
}

func TestPromptTemplate_BuildPrompt_Truncation(t *testing.T) {
	tmpl := &PromptTemplate{
		Category:    "generic",
		Style:       "일반 요약",
		Sections:    []string{"요약"},
		Instruction: "요약하세요.",
	}

	longContent := strings.Repeat("a", 7000)
	prompt := tmpl.BuildPrompt(longContent)

	if !strings.Contains(prompt, "내용이 잘렸습니다") {
		t.Error("long content should be truncated")
	}
}

func TestTruncateContent(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"short", "hello", 10, "hello"},
		{"exact", "hello", 5, "hello"},
		{"truncated", "hello world", 5, "hello\n\n... (내용이 잘렸습니다)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateContent(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncateContent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLoadTemplateFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)

	_, err := loadTemplateFile(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestTemplateContentIntegrity(t *testing.T) {
	dir := findPromptsDir(t)
	reg, err := LoadTemplates(dir)
	if err != nil {
		t.Fatalf("LoadTemplates() error: %v", err)
	}

	// Verify each template has the correct sections per the ticket requirements
	expectedSections := map[model.ContentCategory][]string{
		model.CategoryPrinciple: {"핵심 원리", "작동 방식", "전제조건", "한계점"},
		model.CategoryReview:    {"장점", "단점", "사용 맥락", "추천 대상", "결론"},
		model.CategoryOpinion:   {"주장", "근거", "반론", "저자의 결론"},
		model.CategoryTechIntro: {"핵심 기능", "기존 대비 차이점", "사용 사례", "시작 방법"},
		model.CategoryTutorial:  {"목표", "필요 사전지식", "주요 단계", "핵심 코드/명령어"},
		model.CategoryNews:      {"핵심 사실", "영향", "관련 맥락", "전망"},
	}

	for cat, wantSections := range expectedSections {
		t.Run(string(cat), func(t *testing.T) {
			tmpl := reg.Get(cat)
			if len(tmpl.Sections) != len(wantSections) {
				t.Errorf("sections count = %d, want %d", len(tmpl.Sections), len(wantSections))
			}
			for i, want := range wantSections {
				if i >= len(tmpl.Sections) {
					break
				}
				if tmpl.Sections[i] != want {
					t.Errorf("section[%d] = %q, want %q", i, tmpl.Sections[i], want)
				}
			}

			// Verify the instruction mentions each section
			for _, section := range wantSections {
				if !strings.Contains(tmpl.Instruction, section) {
					t.Errorf("instruction should mention section %q", section)
				}
			}
		})
	}
}

// findPromptsDir locates the prompts directory relative to the test file.
func findPromptsDir(t *testing.T) string {
	t.Helper()
	// Walk up from the test directory to find backend/prompts
	candidates := []string{
		"../../../prompts",
		"../../prompts",
	}
	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		if _, err := os.Stat(abs); err == nil {
			return abs
		}
	}
	t.Skip("prompts directory not found")
	return ""
}
