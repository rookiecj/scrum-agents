package handler

import (
	"encoding/json"
	"log/slog"
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
			slog.Warn("detect: invalid request body",
				slog.String("handler", "detect"),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusBadRequest, DetectResponse{Error: "invalid request body"})
			return
		}

		if req.URL == "" {
			slog.Warn("detect: empty url",
				slog.String("handler", "detect"),
			)
			writeJSON(w, http.StatusBadRequest, DetectResponse{Error: "url is required"})
			return
		}

		linkType, err := urldetect.Detect(req.URL)
		if err != nil {
			slog.Error("detect: invalid URL",
				slog.String("handler", "detect"),
				slog.String("url", req.URL),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusBadRequest, DetectResponse{Error: "invalid URL: " + err.Error()})
			return
		}

		slog.Debug("detect: success",
			slog.String("handler", "detect"),
			slog.String("url", req.URL),
			slog.String("link_type", string(linkType)),
		)
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
		model.LinkTypeArticle:    extractor.NewArticleExtractor(),
		model.LinkTypeYouTube:    extractor.NewYouTubeExtractor(),
		model.LinkTypePDF:        extractor.NewPDFExtractor(),
		model.LinkTypeTwitter:    extractor.NewTwitterExtractor(),
		model.LinkTypeNewsletter: extractor.NewNewsletterExtractor(),
	}

	// Fallback extractor for unsupported types
	fallback := extractor.NewArticleExtractor()

	return func(w http.ResponseWriter, r *http.Request) {
		var req ExtractRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Warn("extract: invalid request body",
				slog.String("handler", "extract"),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusBadRequest, ExtractResponse{Error: "invalid request body"})
			return
		}

		if req.URL == "" {
			slog.Warn("extract: empty url",
				slog.String("handler", "extract"),
			)
			writeJSON(w, http.StatusBadRequest, ExtractResponse{Error: "url is required"})
			return
		}

		linkType, err := urldetect.Detect(req.URL)
		if err != nil {
			slog.Error("extract: invalid URL",
				slog.String("handler", "extract"),
				slog.String("url", req.URL),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusBadRequest, ExtractResponse{Error: "invalid URL: " + err.Error()})
			return
		}

		ext, ok := extractors[linkType]
		if !ok {
			slog.Warn("extract: no extractor for type, using fallback",
				slog.String("handler", "extract"),
				slog.String("url", req.URL),
				slog.String("link_type", string(linkType)),
			)
			// Graceful fallback: attempt generic HTML extraction
			ext = fallback
		}

		result, err := ext.Extract(req.URL)
		if err != nil {
			slog.Error("extract: extraction failed",
				slog.String("handler", "extract"),
				slog.String("url", req.URL),
				slog.String("link_type", string(linkType)),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusInternalServerError, ExtractResponse{Error: "extraction failed: " + err.Error()})
			return
		}

		slog.Debug("extract: success",
			slog.String("handler", "extract"),
			slog.String("url", req.URL),
			slog.String("link_type", string(linkType)),
		)
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
