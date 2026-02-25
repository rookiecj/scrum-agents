package summarizer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// PromptTemplate defines a category-specific summarization prompt.
type PromptTemplate struct {
	Category    string   `json:"category"`
	Style       string   `json:"style"`
	Sections    []string `json:"sections"`
	Instruction string   `json:"instruction"`
}

// TemplateRegistry holds all loaded prompt templates.
type TemplateRegistry struct {
	templates map[model.ContentCategory]*PromptTemplate
	generic   *PromptTemplate
}

// categoryFileMap maps content categories to their template file names.
var categoryFileMap = map[model.ContentCategory]string{
	model.CategoryPrinciple: "principle.json",
	model.CategoryReview:    "review.json",
	model.CategoryOpinion:   "opinion.json",
	model.CategoryTechIntro: "techintro.json",
	model.CategoryTutorial:  "tutorial.json",
	model.CategoryNews:      "news.json",
}

// LoadTemplates loads all prompt templates from the given directory.
// It returns an error if any of the 6 required category templates are missing.
func LoadTemplates(dir string) (*TemplateRegistry, error) {
	reg := &TemplateRegistry{
		templates: make(map[model.ContentCategory]*PromptTemplate),
	}

	// Load all 6 category templates
	for cat, filename := range categoryFileMap {
		tmpl, err := loadTemplateFile(filepath.Join(dir, filename))
		if err != nil {
			return nil, fmt.Errorf("loading template for %s (%s): %w", cat, filename, err)
		}
		reg.templates[cat] = tmpl
	}

	// Load generic fallback template
	generic, err := loadTemplateFile(filepath.Join(dir, "generic.json"))
	if err != nil {
		return nil, fmt.Errorf("loading generic template: %w", err)
	}
	reg.generic = generic

	return reg, nil
}

// Validate checks that all required templates are present and well-formed.
func (r *TemplateRegistry) Validate() error {
	for _, cat := range model.AllCategories() {
		tmpl, ok := r.templates[cat]
		if !ok {
			return fmt.Errorf("missing template for category: %s", cat)
		}
		if tmpl.Instruction == "" {
			return fmt.Errorf("empty instruction in template for category: %s", cat)
		}
		if len(tmpl.Sections) == 0 {
			return fmt.Errorf("no sections defined in template for category: %s", cat)
		}
	}
	if r.generic == nil {
		return fmt.Errorf("missing generic fallback template")
	}
	return nil
}

// Get returns the template for the given category, or the generic template if not found.
func (r *TemplateRegistry) Get(category model.ContentCategory) *PromptTemplate {
	if tmpl, ok := r.templates[category]; ok {
		return tmpl
	}
	return r.generic
}

// GetGeneric returns the generic fallback template.
func (r *TemplateRegistry) GetGeneric() *PromptTemplate {
	return r.generic
}

// Categories returns all registered category names.
func (r *TemplateRegistry) Categories() []model.ContentCategory {
	cats := make([]model.ContentCategory, 0, len(r.templates))
	for cat := range r.templates {
		cats = append(cats, cat)
	}
	return cats
}

// BuildPrompt constructs the full LLM prompt from a template and content.
func (t *PromptTemplate) BuildPrompt(content string) string {
	var sb strings.Builder
	sb.WriteString("당신은 전문 콘텐츠 요약기입니다.\n\n")
	sb.WriteString(fmt.Sprintf("요약 스타일: %s\n\n", t.Style))
	sb.WriteString(t.Instruction)
	sb.WriteString("\n\n---\n\n")
	sb.WriteString(truncateContent(content, 6000))
	sb.WriteString("\n\n---\n\n")
	sb.WriteString("위 글을 한국어로 요약하세요. 마크다운 형식으로 작성하세요.")
	return sb.String()
}

func loadTemplateFile(path string) (*PromptTemplate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	var tmpl PromptTemplate
	if err := json.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("parsing template %s: %w", path, err)
	}

	return &tmpl, nil
}

func truncateContent(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "\n\n... (내용이 잘렸습니다)"
}
