package handler

import (
	"log/slog"
	"net/http"
	"os"
)

// ProviderInfo describes a provider's availability for the frontend.
type ProviderInfo struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
	EnvVar    string `json:"envVar"`
}

// KnownProviders defines the list of known providers with their env var names.
var KnownProviders = []struct {
	Name   string
	EnvVar string
}{
	{Name: "claude", EnvVar: "ANTHROPIC_API_KEY"},
	{Name: "openai", EnvVar: "OPENAI_API_KEY"},
	{Name: "gemini", EnvVar: "GOOGLE_API_KEY"},
}

// HandleProviders returns a handler that lists all known LLM providers with availability status.
// Availability is determined by checking whether the corresponding environment variable is set.
func HandleProviders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providers := make([]ProviderInfo, 0, len(KnownProviders))
		for _, kp := range KnownProviders {
			available := os.Getenv(kp.EnvVar) != ""
			providers = append(providers, ProviderInfo{
				Name:      kp.Name,
				Available: available,
				EnvVar:    kp.EnvVar,
			})
		}

		slog.Debug("providers: listing available providers",
			slog.String("handler", "providers"),
			slog.Int("count", len(providers)),
		)

		writeJSON(w, http.StatusOK, providers)
	}
}
