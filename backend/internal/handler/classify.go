package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rookiecj/scrum-agents/backend/internal/classifier"
	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

type ClassifyRequest struct {
	Content string `json:"content"`
}

type ClassifyResponse struct {
	Classification *model.ClassificationResult `json:"classification,omitempty"`
	Error          string                      `json:"error,omitempty"`
}

// HandleClassify returns a handler that classifies content.
func HandleClassify(cls classifier.Classifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ClassifyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ClassifyResponse{Error: "invalid request body"})
			return
		}

		if req.Content == "" {
			writeJSON(w, http.StatusBadRequest, ClassifyResponse{Error: "content is required"})
			return
		}

		result, err := cls.Classify(req.Content)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ClassifyResponse{Error: "classification failed: " + err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, ClassifyResponse{Classification: result})
	}
}
