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

// OpenAIProvider implements the Provider interface for OpenAI's API.
type OpenAIProvider struct {
	config     Config
	client     *http.Client
	baseURL    string
	classifier *classifier.LLMClassifier
}

// NewOpenAIProvider creates a new OpenAI provider.
func NewOpenAIProvider(config Config) *OpenAIProvider {
	p := &OpenAIProvider{
		config:  config,
		client:  &http.Client{Timeout: config.Timeout},
		baseURL: "https://api.openai.com/v1",
	}
	p.classifier = classifier.NewLLMClassifier(p)
	return p
}

// Name returns the provider type.
func (p *OpenAIProvider) Name() ProviderType {
	return ProviderOpenAI
}

type openaiRequest struct {
	Model     string          `json:"model"`
	Messages  []openaiMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens"`
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Complete sends a prompt to OpenAI and returns the response.
func (p *OpenAIProvider) Complete(prompt string) (string, error) {
	reqBody := openaiRequest{
		Model: p.config.Model,
		Messages: []openaiMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens: p.config.MaxTokens,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", p.baseURL+"/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	var result openaiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenAI")
	}

	return result.Choices[0].Message.Content, nil
}

// Classify classifies content using OpenAI.
func (p *OpenAIProvider) Classify(content string) (*model.ClassificationResult, error) {
	return p.classifier.Classify(content)
}

// Summarize generates a summary using OpenAI.
func (p *OpenAIProvider) Summarize(content string, category model.ContentCategory) (string, error) {
	prompt := fmt.Sprintf("Summarize the following %s content concisely:\n\n%s", string(category), content)
	return p.Complete(prompt)
}
