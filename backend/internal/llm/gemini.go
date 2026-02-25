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

// GeminiProvider implements the Provider interface for Google's Gemini API.
type GeminiProvider struct {
	config     Config
	client     *http.Client
	baseURL    string
	classifier *classifier.LLMClassifier
}

// NewGeminiProvider creates a new Gemini provider.
func NewGeminiProvider(config Config) *GeminiProvider {
	p := &GeminiProvider{
		config:  config,
		client:  &http.Client{Timeout: config.Timeout},
		baseURL: "https://generativelanguage.googleapis.com/v1beta",
	}
	p.classifier = classifier.NewLLMClassifier(p)
	return p
}

// Name returns the provider type.
func (p *GeminiProvider) Name() ProviderType {
	return ProviderGemini
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}

// Complete sends a prompt to Gemini and returns the response.
func (p *GeminiProvider) Complete(prompt string) (string, error) {
	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent", p.baseURL, p.config.Model)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", p.config.APIKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	var result geminiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("Gemini API error: %s", result.Error.Message)
	}

	if len(result.Candidates) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	if len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty content parts from Gemini")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}

// Classify classifies content using Gemini.
func (p *GeminiProvider) Classify(content string) (*model.ClassificationResult, error) {
	return p.classifier.Classify(content)
}

// Summarize generates a summary using Gemini with a pre-built prompt.
func (p *GeminiProvider) Summarize(content string, category model.ContentCategory) (string, error) {
	prompt := fmt.Sprintf("Summarize the following %s content concisely:\n\n%s", string(category), content)
	return p.Complete(prompt)
}
