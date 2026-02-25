package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rookiecj/scrum-agents/backend/internal/handler"
	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// --- Response types mirroring handler package ---

type healthResponse struct {
	Status string `json:"status"`
}

type detectResponse struct {
	LinkInfo model.LinkInfo `json:"link_info"`
	Error    string         `json:"error,omitempty"`
}

type extractResponse struct {
	LinkInfo model.LinkInfo `json:"link_info"`
	Content  string         `json:"content"`
	Error    string         `json:"error,omitempty"`
}

// setupAPIServer creates the same mux as cmd/server/main.go and returns an httptest.Server.
func setupAPIServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})
	mux.HandleFunc("POST /api/detect", handler.HandleDetect())
	mux.HandleFunc("POST /api/extract", handler.HandleExtract())
	return httptest.NewServer(mux)
}

// postJSON sends a POST request with JSON body and returns the response.
func postJSON(t *testing.T, url string, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST %s failed: %v", url, err)
	}
	return resp
}

// decodeJSON decodes the response body into v.
func decodeJSON(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

// --- Health endpoint ---

func TestE2E_Health(t *testing.T) {
	srv := setupAPIServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var body healthResponse
	decodeJSON(t, resp, &body)
	if body.Status != "ok" {
		t.Errorf("status = %q, want %q", body.Status, "ok")
	}
}

// --- Detect endpoint ---

func TestE2E_Detect(t *testing.T) {
	srv := setupAPIServer()
	defer srv.Close()

	tests := []struct {
		name     string
		url      string
		wantType model.LinkType
	}{
		{"article", "https://example.com/blog/post", model.LinkTypeArticle},
		{"youtube", "https://www.youtube.com/watch?v=abc123", model.LinkTypeYouTube},
		{"youtube short", "https://youtu.be/abc123", model.LinkTypeYouTube},
		{"pdf", "https://example.com/paper.pdf", model.LinkTypePDF},
		{"twitter", "https://twitter.com/user/status/123", model.LinkTypeTwitter},
		{"x.com", "https://x.com/user/status/123", model.LinkTypeTwitter},
		{"substack", "https://newsletter.substack.com/p/hello", model.LinkTypeNewsletter},
		{"medium", "https://medium.com/@user/article", model.LinkTypeNewsletter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := postJSON(t, srv.URL+"/api/detect", map[string]string{"url": tt.url})
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status = %d, want 200", resp.StatusCode)
			}

			var dr detectResponse
			decodeJSON(t, resp, &dr)

			if dr.Error != "" {
				t.Errorf("unexpected error: %s", dr.Error)
			}
			if dr.LinkInfo.LinkType != tt.wantType {
				t.Errorf("link_type = %q, want %q", dr.LinkInfo.LinkType, tt.wantType)
			}
			if dr.LinkInfo.URL != tt.url {
				t.Errorf("url = %q, want %q", dr.LinkInfo.URL, tt.url)
			}
		})
	}
}

func TestE2E_Detect_Errors(t *testing.T) {
	srv := setupAPIServer()
	defer srv.Close()

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "empty body",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "url is required",
		},
		{
			name:       "invalid json",
			body:       `not json`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "invalid request body",
		},
		{
			name:       "missing url field",
			body:       `{"link":"https://example.com"}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "url is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Post(srv.URL+"/api/detect", "application/json", strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("POST failed: %v", err)
			}
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			var dr detectResponse
			decodeJSON(t, resp, &dr)
			if !strings.Contains(dr.Error, tt.wantErr) {
				t.Errorf("error = %q, want containing %q", dr.Error, tt.wantErr)
			}
		})
	}
}

// --- Extract endpoint ---

func TestE2E_Extract_Article(t *testing.T) {
	// Mock external HTML server
	externalSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html>
<head><title>E2E Test Article</title></head>
<body>
<p>This is the main content of the end-to-end test article.</p>
<p>It has multiple paragraphs to verify extraction works correctly.</p>
</body>
</html>`))
	}))
	defer externalSrv.Close()

	srv := setupAPIServer()
	defer srv.Close()

	// Step 1: Detect — localhost URL is classified as "article"
	resp := postJSON(t, srv.URL+"/api/detect", map[string]string{"url": externalSrv.URL + "/blog/post"})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("detect status = %d, want 200", resp.StatusCode)
	}
	var dr detectResponse
	decodeJSON(t, resp, &dr)
	if dr.LinkInfo.LinkType != model.LinkTypeArticle {
		t.Errorf("detect link_type = %q, want %q", dr.LinkInfo.LinkType, model.LinkTypeArticle)
	}

	// Step 2: Extract — article content extracted from mock server
	resp = postJSON(t, srv.URL+"/api/extract", map[string]string{"url": externalSrv.URL + "/blog/post"})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("extract status = %d, want 200", resp.StatusCode)
	}
	var er extractResponse
	decodeJSON(t, resp, &er)

	if er.Error != "" {
		t.Errorf("unexpected error: %s", er.Error)
	}
	if er.LinkInfo.Title != "E2E Test Article" {
		t.Errorf("title = %q, want %q", er.LinkInfo.Title, "E2E Test Article")
	}
	if er.LinkInfo.LinkType != model.LinkTypeArticle {
		t.Errorf("link_type = %q, want %q", er.LinkInfo.LinkType, model.LinkTypeArticle)
	}
	if !strings.Contains(er.Content, "main content") {
		t.Errorf("content should contain 'main content', got %q", er.Content)
	}
	if !strings.Contains(er.Content, "multiple paragraphs") {
		t.Errorf("content should contain 'multiple paragraphs', got %q", er.Content)
	}
}

func TestE2E_Extract_PDF(t *testing.T) {
	// Mock server that serves a valid PDF at a .pdf path
	simplePDF := buildSimplePDF("Hello from E2E PDF test")

	externalSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.Write(simplePDF)
	}))
	defer externalSrv.Close()

	srv := setupAPIServer()
	defer srv.Close()

	// Use a .pdf path so urldetect classifies it as PDF
	pdfURL := externalSrv.URL + "/paper.pdf"

	// Step 1: Detect
	resp := postJSON(t, srv.URL+"/api/detect", map[string]string{"url": pdfURL})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("detect status = %d, want 200", resp.StatusCode)
	}
	var dr detectResponse
	decodeJSON(t, resp, &dr)
	if dr.LinkInfo.LinkType != model.LinkTypePDF {
		t.Errorf("detect link_type = %q, want %q", dr.LinkInfo.LinkType, model.LinkTypePDF)
	}

	// Step 2: Extract
	resp = postJSON(t, srv.URL+"/api/extract", map[string]string{"url": pdfURL})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("extract status = %d, want 200", resp.StatusCode)
	}
	var er extractResponse
	decodeJSON(t, resp, &er)

	if er.Error != "" {
		t.Errorf("unexpected error: %s", er.Error)
	}
	if er.LinkInfo.LinkType != model.LinkTypePDF {
		t.Errorf("link_type = %q, want %q", er.LinkInfo.LinkType, model.LinkTypePDF)
	}
	if !strings.Contains(er.Content, "Hello from E2E PDF test") {
		t.Errorf("content should contain PDF text, got %q", er.Content)
	}
}

func TestE2E_Extract_Fallback(t *testing.T) {
	// Mock server for a URL that would be "newsletter" type but serves HTML
	// Since localhost won't match newsletter hostnames, this tests the article fallback path
	externalSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head><title>Fallback Article</title></head>
<body><p>Content extracted via fallback path.</p></body></html>`))
	}))
	defer externalSrv.Close()

	srv := setupAPIServer()
	defer srv.Close()

	resp := postJSON(t, srv.URL+"/api/extract", map[string]string{"url": externalSrv.URL + "/post"})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("extract status = %d, want 200", resp.StatusCode)
	}
	var er extractResponse
	decodeJSON(t, resp, &er)

	if er.Error != "" {
		t.Errorf("unexpected error: %s", er.Error)
	}
	if er.LinkInfo.Title != "Fallback Article" {
		t.Errorf("title = %q, want %q", er.LinkInfo.Title, "Fallback Article")
	}
}

func TestE2E_Extract_Errors(t *testing.T) {
	srv := setupAPIServer()
	defer srv.Close()

	t.Run("empty body", func(t *testing.T) {
		resp, err := http.Post(srv.URL+"/api/extract", "application/json", strings.NewReader(`{}`))
		if err != nil {
			t.Fatalf("POST failed: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", resp.StatusCode)
		}
		var er extractResponse
		decodeJSON(t, resp, &er)
		if !strings.Contains(er.Error, "url is required") {
			t.Errorf("error = %q, want containing 'url is required'", er.Error)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		resp, err := http.Post(srv.URL+"/api/extract", "application/json", strings.NewReader(`not json`))
		if err != nil {
			t.Fatalf("POST failed: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", resp.StatusCode)
		}
		var er extractResponse
		decodeJSON(t, resp, &er)
		if !strings.Contains(er.Error, "invalid request body") {
			t.Errorf("error = %q, want containing 'invalid request body'", er.Error)
		}
	})

	t.Run("external server error", func(t *testing.T) {
		errorSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer errorSrv.Close()

		resp := postJSON(t, srv.URL+"/api/extract", map[string]string{"url": errorSrv.URL + "/fail"})
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", resp.StatusCode)
		}
		var er extractResponse
		decodeJSON(t, resp, &er)
		if !strings.Contains(er.Error, "extraction failed") {
			t.Errorf("error = %q, want containing 'extraction failed'", er.Error)
		}
	})

	t.Run("pdf too large", func(t *testing.T) {
		largePDFSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", 11*1024*1024)) // 11MB
			w.Write([]byte("%PDF-1.4"))
		}))
		defer largePDFSrv.Close()

		resp := postJSON(t, srv.URL+"/api/extract", map[string]string{"url": largePDFSrv.URL + "/large.pdf"})
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", resp.StatusCode)
		}
		var er extractResponse
		decodeJSON(t, resp, &er)
		if !strings.Contains(er.Error, "maximum size") {
			t.Errorf("error = %q, want containing 'maximum size'", er.Error)
		}
	})
}

// --- Full pipeline: detect → extract ---

func TestE2E_DetectThenExtract_Pipeline(t *testing.T) {
	// This test simulates the full user workflow: detect the type, then extract content
	externalSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, ".pdf"):
			w.Header().Set("Content-Type", "application/pdf")
			w.Write(buildSimplePDF("Pipeline PDF content"))
		default:
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<html><head><title>Pipeline Article</title></head>
<body><p>Pipeline article content for E2E verification.</p></body></html>`))
		}
	}))
	defer externalSrv.Close()

	srv := setupAPIServer()
	defer srv.Close()

	tests := []struct {
		name         string
		path         string
		wantType     model.LinkType
		wantTitle    string
		wantContent  string
	}{
		{
			name:        "article pipeline",
			path:        "/blog/my-post",
			wantType:    model.LinkTypeArticle,
			wantTitle:   "Pipeline Article",
			wantContent: "Pipeline article content",
		},
		{
			name:        "pdf pipeline",
			path:        "/docs/paper.pdf",
			wantType:    model.LinkTypePDF,
			wantContent: "Pipeline PDF content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetURL := externalSrv.URL + tt.path

			// Phase 1: Detect
			detectResp := postJSON(t, srv.URL+"/api/detect", map[string]string{"url": targetURL})
			if detectResp.StatusCode != http.StatusOK {
				t.Fatalf("detect: status = %d, want 200", detectResp.StatusCode)
			}
			var dr detectResponse
			decodeJSON(t, detectResp, &dr)

			if dr.LinkInfo.LinkType != tt.wantType {
				t.Errorf("detect: link_type = %q, want %q", dr.LinkInfo.LinkType, tt.wantType)
			}

			// Phase 2: Extract
			extractResp := postJSON(t, srv.URL+"/api/extract", map[string]string{"url": targetURL})
			if extractResp.StatusCode != http.StatusOK {
				t.Fatalf("extract: status = %d, want 200", extractResp.StatusCode)
			}
			var er extractResponse
			decodeJSON(t, extractResp, &er)

			if er.Error != "" {
				t.Errorf("extract: unexpected error: %s", er.Error)
			}
			if er.LinkInfo.LinkType != tt.wantType {
				t.Errorf("extract: link_type = %q, want %q", er.LinkInfo.LinkType, tt.wantType)
			}
			if tt.wantTitle != "" && er.LinkInfo.Title != tt.wantTitle {
				t.Errorf("extract: title = %q, want %q", er.LinkInfo.Title, tt.wantTitle)
			}
			if tt.wantContent != "" && !strings.Contains(er.Content, tt.wantContent) {
				t.Errorf("extract: content should contain %q, got %q", tt.wantContent, er.Content)
			}
		})
	}
}

// buildSimplePDF creates a minimal valid PDF with extractable text.
func buildSimplePDF(text string) []byte {
	return []byte(fmt.Sprintf(`%%PDF-1.4
1 0 obj <</Type /Catalog /Pages 2 0 R>> endobj
2 0 obj <</Type /Pages /Kids [3 0 R] /Count 1>> endobj
3 0 obj <</Type /Page /Parent 2 0 R /MediaBox [0 0 612 792]
/Contents 4 0 R /Resources <</Font <</F1 5 0 R>>>>>> endobj
4 0 obj <</Length 44>>
stream
BT /F1 12 Tf 100 700 Td (%s) Tj ET
endstream endobj
5 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica>> endobj
xref
0 6
trailer <</Size 6 /Root 1 0 R>>
startxref
0
%%%%EOF`, text))
}
