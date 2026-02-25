package summarizer

import (
	"fmt"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// LLMClient is the interface for making LLM API calls.
type LLMClient interface {
	Complete(prompt string) (string, error)
}

// SummaryResult holds the summarization result with metadata.
type SummaryResult struct {
	Summary        string                `json:"summary"`
	Category       model.ContentCategory `json:"category"`
	Style          string                `json:"style"`
	LowConfidence  bool                  `json:"low_confidence,omitempty"`
	TemplateUsed   string                `json:"template_used"`
}

// Summarizer generates category-optimized summaries using prompt templates.
type Summarizer struct {
	registry            *TemplateRegistry
	confidenceThreshold float64
}

// NewSummarizer creates a Summarizer with the given template registry.
// confidenceThreshold determines when to fall back to the generic template.
func NewSummarizer(registry *TemplateRegistry, confidenceThreshold float64) *Summarizer {
	if confidenceThreshold <= 0 {
		confidenceThreshold = 0.6
	}
	return &Summarizer{
		registry:            registry,
		confidenceThreshold: confidenceThreshold,
	}
}

// Summarize generates a summary using the appropriate template based on classification.
func (s *Summarizer) Summarize(client LLMClient, content string, classification *model.ClassificationResult) (*SummaryResult, error) {
	var tmpl *PromptTemplate
	var lowConfidence bool

	if classification.Confidence < s.confidenceThreshold {
		tmpl = s.registry.GetGeneric()
		lowConfidence = true
	} else {
		tmpl = s.registry.Get(classification.Primary)
	}

	prompt := tmpl.BuildPrompt(content)
	summary, err := client.Complete(prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM summarization failed: %w", err)
	}

	return &SummaryResult{
		Summary:       summary,
		Category:      classification.Primary,
		Style:         tmpl.Style,
		LowConfidence: lowConfidence,
		TemplateUsed:  tmpl.Category,
	}, nil
}

// SummarizeWithCategory generates a summary using the template for the given category directly.
func (s *Summarizer) SummarizeWithCategory(client LLMClient, content string, category model.ContentCategory) (*SummaryResult, error) {
	tmpl := s.registry.Get(category)

	prompt := tmpl.BuildPrompt(content)
	summary, err := client.Complete(prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM summarization failed: %w", err)
	}

	return &SummaryResult{
		Summary:      summary,
		Category:     category,
		Style:        tmpl.Style,
		TemplateUsed: tmpl.Category,
	}, nil
}

// Registry returns the template registry for inspection.
func (s *Summarizer) Registry() *TemplateRegistry {
	return s.registry
}
