package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rookiecj/scrum-agents/backend/internal/extractor"
	"github.com/rookiecj/scrum-agents/backend/internal/model"
	"github.com/rookiecj/scrum-agents/backend/internal/urldetect"
)

type DetectRequest struct {
	URL string `json:"url"`
}

type DetectResponse struct {
	LinkInfo model.LinkInfo `json:"link_info"`
	Error    string         `json:"error,omitempty"`
}

type ExtractRequest struct {
	URL string `json:"url"`
}

type ExtractResponse struct {
	LinkInfo model.LinkInfo `json:"link_info"`
	Content  string         `json:"content"`
	Error    string         `json:"error,omitempty"`
}

// HandleDetect returns a handler that detects the type of a URL.
func HandleDetect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DetectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, DetectResponse{Error: "invalid request body"})
			return
		}

		if req.URL == "" {
			writeJSON(w, http.StatusBadRequest, DetectResponse{Error: "url is required"})
			return
		}

		linkType, err := urldetect.Detect(req.URL)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, DetectResponse{Error: "invalid URL: " + err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, DetectResponse{
			LinkInfo: model.LinkInfo{
				URL:      req.URL,
				LinkType: linkType,
			},
		})
	}
}

// HandleExtract returns a handler that extracts content from a URL.
func HandleExtract() http.HandlerFunc {
	extractors := map[model.LinkType]extractor.Extractor{
		model.LinkTypeArticle: extractor.NewArticleExtractor(),
		model.LinkTypeYouTube: extractor.NewYouTubeExtractor(),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req ExtractRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ExtractResponse{Error: "invalid request body"})
			return
		}

		if req.URL == "" {
			writeJSON(w, http.StatusBadRequest, ExtractResponse{Error: "url is required"})
			return
		}

		linkType, err := urldetect.Detect(req.URL)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ExtractResponse{Error: "invalid URL: " + err.Error()})
			return
		}

		ext, ok := extractors[linkType]
		if !ok {
			writeJSON(w, http.StatusOK, ExtractResponse{
				LinkInfo: model.LinkInfo{URL: req.URL, LinkType: linkType},
				Error:    "extraction not yet supported for type: " + string(linkType),
			})
			return
		}

		result, err := ext.Extract(req.URL)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ExtractResponse{Error: "extraction failed: " + err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, ExtractResponse{
			LinkInfo: result.LinkInfo,
			Content:  result.Content,
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
