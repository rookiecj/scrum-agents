package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rookiecj/scrum-agents/backend/internal/classifier"
	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// ClaudeProvider implements the Provider interface for Anthropic's Claude API.
type ClaudeProvider struct {
	config     Config
	client     *http.Client
	baseURL    string
	classifier *classifier.LLMClassifier
}

// NewClaudeProvider creates a new Claude provider.
func NewClaudeProvider(config Config) *ClaudeProvider {
	p := &ClaudeProvider{
		config:  config,
		client:  &http.Client{Timeout: config.Timeout},
		baseURL: "https://api.anthropic.com/v1",
	}
	p.classifier = classifier.NewLLMClassifier(p)
	return p
}

// Name returns the provider type.
func (p *ClaudeProvider) Name() ProviderType {
	return ProviderClaude
}

type claudeRequest struct {
	Model     string           `json:"model"`
	MaxTokens int              `json:"max_tokens"`
	Messages  []claudeMessage  `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Complete sends a prompt to Claude and returns the response.
func (p *ClaudeProvider) Complete(prompt string) (string, error) {
	reqBody := claudeRequest{
		Model:     p.config.Model,
		MaxTokens: p.config.MaxTokens,
		Messages: []claudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", p.baseURL+"/messages", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling Claude API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	var result claudeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("Claude API error: %s", result.Error.Message)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return result.Content[0].Text, nil
}

// Classify classifies content using Claude.
func (p *ClaudeProvider) Classify(content string) (*model.ClassificationResult, error) {
	return p.classifier.Classify(content)
}

// Summarize generates a summary using Claude.
func (p *ClaudeProvider) Summarize(content string, category model.ContentCategory) (string, error) {
	prompt := fmt.Sprintf("Summarize the following %s content concisely:\n\n%s", string(category), content)
	return p.Complete(prompt)
}
