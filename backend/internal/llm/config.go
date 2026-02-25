package llm

import "time"

// Config holds configuration for an LLM provider.
type Config struct {
	APIKey     string        `json:"api_key"`
	Model      string        `json:"model"`
	MaxTokens  int           `json:"max_tokens"`
	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"max_retries"`
}

// DefaultClaudeConfig returns default configuration for Claude.
func DefaultClaudeConfig(apiKey string) Config {
	return Config{
		APIKey:     apiKey,
		Model:      "claude-sonnet-4-6",
		MaxTokens:  4096,
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}
}

// DefaultOpenAIConfig returns default configuration for OpenAI.
func DefaultOpenAIConfig(apiKey string) Config {
	return Config{
		APIKey:     apiKey,
		Model:      "gpt-4o",
		MaxTokens:  4096,
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}
}

// DefaultGeminiConfig returns default configuration for Google Gemini.
func DefaultGeminiConfig(apiKey string) Config {
	return Config{
		APIKey:     apiKey,
		Model:      "gemini-2.0-flash",
		MaxTokens:  4096,
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}
}
