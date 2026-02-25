package llm

import (
	"fmt"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// Adapter manages multiple LLM providers and routes requests to the selected one.
type Adapter struct {
	providers      map[ProviderType]Provider
	defaultProvider ProviderType
}

// NewAdapter creates an Adapter with the given providers and default.
func NewAdapter(defaultProvider ProviderType, providers ...Provider) (*Adapter, error) {
	if len(providers) == 0 {
		return nil, fmt.Errorf("at least one provider is required")
	}

	a := &Adapter{
		providers:      make(map[ProviderType]Provider),
		defaultProvider: defaultProvider,
	}

	for _, p := range providers {
		a.providers[p.Name()] = p
	}

	if _, ok := a.providers[defaultProvider]; !ok {
		return nil, fmt.Errorf("default provider %q not found in registered providers", defaultProvider)
	}

	return a, nil
}

// GetProvider returns the provider for the given type, or the default if empty.
func (a *Adapter) GetProvider(providerType ProviderType) (Provider, error) {
	if providerType == "" {
		providerType = a.defaultProvider
	}

	p, ok := a.providers[providerType]
	if !ok {
		return nil, fmt.Errorf("provider %q not available", providerType)
	}
	return p, nil
}

// Complete sends a prompt using the specified or default provider.
func (a *Adapter) Complete(prompt string, providerType ProviderType) (string, error) {
	p, err := a.GetProvider(providerType)
	if err != nil {
		return "", err
	}
	return p.Complete(prompt)
}

// Classify classifies content using the specified or default provider.
func (a *Adapter) Classify(content string, providerType ProviderType) (*model.ClassificationResult, error) {
	p, err := a.GetProvider(providerType)
	if err != nil {
		return nil, err
	}
	return p.Classify(content)
}

// Summarize generates a summary using the specified or default provider.
func (a *Adapter) Summarize(content string, category model.ContentCategory, providerType ProviderType) (string, error) {
	p, err := a.GetProvider(providerType)
	if err != nil {
		return "", err
	}
	return p.Summarize(content, category)
}

// AvailableProviders returns the list of registered provider types.
func (a *Adapter) AvailableProviders() []ProviderType {
	var types []ProviderType
	for t := range a.providers {
		types = append(types, t)
	}
	return types
}

// DefaultProvider returns the default provider type.
func (a *Adapter) DefaultProvider() ProviderType {
	return a.defaultProvider
}
