package llm

import (
	"fmt"
	"testing"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// mockProvider is a test double for Provider.
type mockProvider struct {
	name        ProviderType
	completeRes string
	completeErr error
	classifyRes *model.ClassificationResult
	classifyErr error
	summarizeRes string
	summarizeErr error
}

func (m *mockProvider) Complete(prompt string) (string, error) {
	return m.completeRes, m.completeErr
}

func (m *mockProvider) Classify(content string) (*model.ClassificationResult, error) {
	return m.classifyRes, m.classifyErr
}

func (m *mockProvider) Summarize(content string, category model.ContentCategory) (string, error) {
	return m.summarizeRes, m.summarizeErr
}

func (m *mockProvider) Name() ProviderType {
	return m.name
}

func TestNewAdapter(t *testing.T) {
	tests := []struct {
		name           string
		defaultProvider ProviderType
		providers      []Provider
		wantErr        bool
	}{
		{
			name:           "single provider",
			defaultProvider: ProviderClaude,
			providers:      []Provider{&mockProvider{name: ProviderClaude}},
		},
		{
			name:           "multiple providers",
			defaultProvider: ProviderClaude,
			providers: []Provider{
				&mockProvider{name: ProviderClaude},
				&mockProvider{name: ProviderOpenAI},
			},
		},
		{
			name:           "no providers",
			defaultProvider: ProviderClaude,
			providers:      []Provider{},
			wantErr:        true,
		},
		{
			name:           "default not found",
			defaultProvider: ProviderClaude,
			providers:      []Provider{&mockProvider{name: ProviderOpenAI}},
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewAdapter(tt.defaultProvider, tt.providers...)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if a == nil {
				t.Fatal("expected non-nil adapter")
			}
		})
	}
}

func TestAdapter_GetProvider(t *testing.T) {
	claude := &mockProvider{name: ProviderClaude}
	openai := &mockProvider{name: ProviderOpenAI}
	adapter, _ := NewAdapter(ProviderClaude, claude, openai)

	tests := []struct {
		name     string
		provider ProviderType
		want     ProviderType
		wantErr  bool
	}{
		{name: "get claude", provider: ProviderClaude, want: ProviderClaude},
		{name: "get openai", provider: ProviderOpenAI, want: ProviderOpenAI},
		{name: "get default (empty)", provider: "", want: ProviderClaude},
		{name: "unknown provider", provider: "unknown", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := adapter.GetProvider(tt.provider)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.Name() != tt.want {
				t.Errorf("provider = %q, want %q", p.Name(), tt.want)
			}
		})
	}
}

func TestAdapter_Complete(t *testing.T) {
	claude := &mockProvider{name: ProviderClaude, completeRes: "claude response"}
	openai := &mockProvider{name: ProviderOpenAI, completeRes: "openai response"}
	adapter, _ := NewAdapter(ProviderClaude, claude, openai)

	tests := []struct {
		name     string
		provider ProviderType
		want     string
		wantErr  bool
	}{
		{name: "claude", provider: ProviderClaude, want: "claude response"},
		{name: "openai", provider: ProviderOpenAI, want: "openai response"},
		{name: "default", provider: "", want: "claude response"},
		{name: "unknown", provider: "unknown", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := adapter.Complete("prompt", tt.provider)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Complete() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAdapter_Classify(t *testing.T) {
	result := &model.ClassificationResult{
		Primary:    model.CategoryTutorial,
		Confidence: 0.95,
	}
	claude := &mockProvider{name: ProviderClaude, classifyRes: result}
	adapter, _ := NewAdapter(ProviderClaude, claude)

	got, err := adapter.Classify("content", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Primary != model.CategoryTutorial {
		t.Errorf("primary = %q, want %q", got.Primary, model.CategoryTutorial)
	}
}

func TestAdapter_Summarize(t *testing.T) {
	claude := &mockProvider{name: ProviderClaude, summarizeRes: "summary text"}
	adapter, _ := NewAdapter(ProviderClaude, claude)

	got, err := adapter.Summarize("content", model.CategoryPrinciple, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "summary text" {
		t.Errorf("Summarize() = %q, want %q", got, "summary text")
	}
}

func TestAdapter_ErrorPropagation(t *testing.T) {
	claude := &mockProvider{
		name:        ProviderClaude,
		completeErr: fmt.Errorf("API down"),
		classifyErr: fmt.Errorf("classify error"),
		summarizeErr: fmt.Errorf("summarize error"),
	}
	adapter, _ := NewAdapter(ProviderClaude, claude)

	_, err := adapter.Complete("prompt", "")
	if err == nil {
		t.Error("expected error from Complete")
	}

	_, err = adapter.Classify("content", "")
	if err == nil {
		t.Error("expected error from Classify")
	}

	_, err = adapter.Summarize("content", model.CategoryNews, "")
	if err == nil {
		t.Error("expected error from Summarize")
	}
}

func TestAdapter_AvailableProviders(t *testing.T) {
	claude := &mockProvider{name: ProviderClaude}
	openai := &mockProvider{name: ProviderOpenAI}
	adapter, _ := NewAdapter(ProviderClaude, claude, openai)

	providers := adapter.AvailableProviders()
	if len(providers) != 2 {
		t.Errorf("got %d providers, want 2", len(providers))
	}
}

func TestAdapter_DefaultProvider(t *testing.T) {
	claude := &mockProvider{name: ProviderClaude}
	adapter, _ := NewAdapter(ProviderClaude, claude)

	if adapter.DefaultProvider() != ProviderClaude {
		t.Errorf("default = %q, want %q", adapter.DefaultProvider(), ProviderClaude)
	}
}
