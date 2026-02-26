package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/rookiecj/scrum-agents/backend/internal/llm"
	"github.com/rookiecj/scrum-agents/backend/internal/model"
	"github.com/rookiecj/scrum-agents/backend/internal/summarizer"
)

// SummarizeRequest is the request body for the summarize endpoint.
// Accepts either a full Classification object or a Category string.
type SummarizeRequest struct {
	Content        string                     `json:"content"`
	Classification *model.ClassificationResult `json:"classification,omitempty"`
	Category       string                     `json:"category,omitempty"`
	Provider       string                     `json:"provider,omitempty"`
}

// SummarizeResponse is the response body for the summarize endpoint.
type SummarizeResponse struct {
	Result *summarizer.SummaryResult `json:"result,omitempty"`
	Error  string                   `json:"error,omitempty"`
}

// HandleSummarize returns a handler that summarizes content using type-specific templates.
// It accepts an optional "provider" field and supports both "classification" (object) and "category" (string).
func HandleSummarize(s *summarizer.Summarizer, defaultClient summarizer.LLMClient, providers map[string]llm.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SummarizeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Warn("summarize: invalid request body",
				slog.String("handler", "summarize"),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusBadRequest, SummarizeResponse{Error: "invalid request body"})
			return
		}

		if req.Content == "" {
			slog.Warn("summarize: empty content",
				slog.String("handler", "summarize"),
			)
			writeJSON(w, http.StatusBadRequest, SummarizeResponse{Error: "content is required"})
			return
		}

		// Build classification from either the full object or the category string
		classification := req.Classification
		if classification == nil && req.Category != "" {
			classification = &model.ClassificationResult{
				Primary:    model.ContentCategory(req.Category),
				Confidence: 1.0,
			}
		}
		if classification == nil {
			slog.Warn("summarize: missing classification and category",
				slog.String("handler", "summarize"),
			)
			writeJSON(w, http.StatusBadRequest, SummarizeResponse{Error: "classification or category is required"})
			return
		}

		// Select LLM client based on requested provider
		var client summarizer.LLMClient = defaultClient
		if req.Provider != "" {
			if p, ok := providers[req.Provider]; ok {
				client = p
			} else {
				writeJSON(w, http.StatusBadRequest, SummarizeResponse{Error: "provider not available: " + req.Provider})
				return
			}
		}

		result, err := s.Summarize(client, req.Content, classification)
		if err != nil {
			slog.Error("summarize: summarization failed",
				slog.String("handler", "summarize"),
				slog.String("category", string(classification.Primary)),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusInternalServerError, SummarizeResponse{Error: "summarization failed: " + err.Error()})
			return
		}

		slog.Debug("summarize: success",
			slog.String("handler", "summarize"),
			slog.String("template_used", result.TemplateUsed),
		)
		writeJSON(w, http.StatusOK, SummarizeResponse{Result: result})
	}
}
