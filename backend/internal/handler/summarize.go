package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
	"github.com/rookiecj/scrum-agents/backend/internal/summarizer"
)

// SummarizeRequest is the request body for the summarize endpoint.
type SummarizeRequest struct {
	Content        string                     `json:"content"`
	Classification *model.ClassificationResult `json:"classification"`
}

// SummarizeResponse is the response body for the summarize endpoint.
type SummarizeResponse struct {
	Result *summarizer.SummaryResult `json:"result,omitempty"`
	Error  string                   `json:"error,omitempty"`
}

// HandleSummarize returns a handler that summarizes content using type-specific templates.
func HandleSummarize(s *summarizer.Summarizer, client summarizer.LLMClient) http.HandlerFunc {
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

		if req.Classification == nil {
			slog.Warn("summarize: missing classification",
				slog.String("handler", "summarize"),
			)
			writeJSON(w, http.StatusBadRequest, SummarizeResponse{Error: "classification is required"})
			return
		}

		result, err := s.Summarize(client, req.Content, req.Classification)
		if err != nil {
			slog.Error("summarize: summarization failed",
				slog.String("handler", "summarize"),
				slog.String("category", string(req.Classification.Primary)),
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
