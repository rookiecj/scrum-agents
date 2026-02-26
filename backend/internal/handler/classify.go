package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/rookiecj/scrum-agents/backend/internal/classifier"
	"github.com/rookiecj/scrum-agents/backend/internal/llm"
	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

type ClassifyRequest struct {
	Content  string `json:"content"`
	Provider string `json:"provider,omitempty"`
}

type ClassifyResponse struct {
	Classification *model.ClassificationResult `json:"classification,omitempty"`
	Error          string                      `json:"error,omitempty"`
}

// HandleClassify returns a handler that classifies content.
// It accepts an optional "provider" field in the request to select the LLM provider.
func HandleClassify(defaultClassifier classifier.Classifier, providers map[string]llm.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ClassifyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Warn("classify: invalid request body",
				slog.String("handler", "classify"),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusBadRequest, ClassifyResponse{Error: "invalid request body"})
			return
		}

		if req.Content == "" {
			slog.Warn("classify: empty content",
				slog.String("handler", "classify"),
			)
			writeJSON(w, http.StatusBadRequest, ClassifyResponse{Error: "content is required"})
			return
		}

		// Select classifier based on requested provider
		var cls classifier.Classifier = defaultClassifier
		if req.Provider != "" {
			if p, ok := providers[req.Provider]; ok {
				cls = p
			} else {
				writeJSON(w, http.StatusBadRequest, ClassifyResponse{Error: "provider not available: " + req.Provider})
				return
			}
		}

		result, err := cls.Classify(req.Content)
		if err != nil {
			slog.Error("classify: classification failed",
				slog.String("handler", "classify"),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusInternalServerError, ClassifyResponse{Error: "classification failed: " + err.Error()})
			return
		}

		slog.Debug("classify: success",
			slog.String("handler", "classify"),
			slog.String("primary", string(result.Primary)),
			slog.Float64("confidence", result.Confidence),
		)
		writeJSON(w, http.StatusOK, ClassifyResponse{Classification: result})
	}
}
