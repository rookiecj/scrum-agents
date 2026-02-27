package classifier

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// LLMClient is the interface for making LLM API calls.
type LLMClient interface {
	Complete(prompt string) (string, error)
}

// LLMClassifier classifies content using an LLM provider.
type LLMClassifier struct {
	Client LLMClient
}

// NewLLMClassifier creates a new LLMClassifier with the given LLM client.
func NewLLMClassifier(client LLMClient) *LLMClassifier {
	return &LLMClassifier{Client: client}
}

// Classify sends the content to the LLM and parses the classification result.
func (c *LLMClassifier) Classify(content string) (*model.ClassificationResult, error) {
	prompt := ClassificationPrompt(content)

	response, err := c.Client.Complete(prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM classification failed: %w", err)
	}

	var result model.ClassificationResult
	cleaned := stripCodeFences(response)
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("parsing classification response: %w", err)
	}

	if !isValidCategory(result.Primary) {
		return nil, fmt.Errorf("invalid primary category: %s", result.Primary)
	}

	return &result, nil
}

func isValidCategory(cat model.ContentCategory) bool {
	for _, valid := range model.AllCategories() {
		if cat == valid {
			return true
		}
	}
	return false
}

// stripCodeFences removes markdown code fences from LLM responses.
// Some providers (e.g., OpenAI) wrap JSON output in ```json ... ```.
func stripCodeFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		// Remove opening fence (e.g., ```json or ```)
		if idx := strings.Index(s, "\n"); idx >= 0 {
			s = s[idx+1:]
		}
		// Remove closing fence
		if idx := strings.LastIndex(s, "```"); idx >= 0 {
			s = s[:idx]
		}
		s = strings.TrimSpace(s)
	}
	return s
}
