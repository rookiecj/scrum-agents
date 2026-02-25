package llm

import (
	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// ProviderType identifies an LLM provider.
type ProviderType string

const (
	ProviderClaude ProviderType = "claude"
	ProviderOpenAI ProviderType = "openai"
	ProviderGemini ProviderType = "gemini"
)

// Provider defines the interface for LLM providers.
type Provider interface {
	// Complete sends a prompt and returns the response text.
	Complete(prompt string) (string, error)

	// Classify classifies content into a category.
	Classify(content string) (*model.ClassificationResult, error)

	// Summarize generates a summary of the content.
	Summarize(content string, category model.ContentCategory) (string, error)

	// Name returns the provider name.
	Name() ProviderType
}
