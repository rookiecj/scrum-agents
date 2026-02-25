package handler

import (
	"encoding/json"
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
			writeJSON(w, http.StatusBadRequest, SummarizeResponse{Error: "invalid request body"})
			return
		}

		if req.Content == "" {
			writeJSON(w, http.StatusBadRequest, SummarizeResponse{Error: "content is required"})
			return
		}

		if req.Classification == nil {
			writeJSON(w, http.StatusBadRequest, SummarizeResponse{Error: "classification is required"})
			return
		}

		result, err := s.Summarize(client, req.Content, req.Classification)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, SummarizeResponse{Error: "summarization failed: " + err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, SummarizeResponse{Result: result})
	}
}
