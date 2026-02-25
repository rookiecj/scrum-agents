package extractor

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

const maxPDFSize = 10 * 1024 * 1024 // 10MB

// PDFExtractor extracts text content from PDF URLs.
type PDFExtractor struct {
	Client *http.Client
}

// NewPDFExtractor creates a new PDFExtractor.
func NewPDFExtractor() *PDFExtractor {
	return &PDFExtractor{
		Client: &http.Client{},
	}
}

// Extract downloads a PDF and extracts text content.
func (e *PDFExtractor) Extract(rawURL string) (*model.ExtractedContent, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LinkSummarizer/1.0)")

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching PDF %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for PDF %s", resp.StatusCode, rawURL)
	}

	// Check Content-Length if available
	if resp.ContentLength > maxPDFSize {
		return nil, fmt.Errorf("PDF exceeds maximum size of 10MB (size: %d bytes)", resp.ContentLength)
	}

	// Read with a limit to enforce size restriction
	limited := io.LimitReader(resp.Body, maxPDFSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("reading PDF body: %w", err)
	}

	if len(data) > maxPDFSize {
		return nil, fmt.Errorf("PDF exceeds maximum size of 10MB (size: >%d bytes)", maxPDFSize)
	}

	title := extractPDFTitle(data)
	text := extractPDFText(data)

	if text == "" {
		return nil, fmt.Errorf("could not extract text from PDF (possibly image-based or encrypted)")
	}

	return &model.ExtractedContent{
		LinkInfo: model.LinkInfo{
			URL:      rawURL,
			LinkType: model.LinkTypePDF,
			Title:    title,
		},
		Content: text,
	}, nil
}

// extractPDFTitle attempts to extract the title from PDF metadata.
func extractPDFTitle(data []byte) string {
	s := string(data)
	// Look for /Title in the PDF info dictionary
	re := regexp.MustCompile(`/Title\s*\(([^)]+)\)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// extractPDFText extracts text from PDF content streams.
// This is a basic extractor that handles simple text PDFs.
func extractPDFText(data []byte) string {
	s := string(data)
	var texts []string

	// Extract text from BT...ET blocks (PDF text objects)
	btRe := regexp.MustCompile(`BT\s([\s\S]*?)ET`)
	blocks := btRe.FindAllStringSubmatch(s, -1)

	for _, block := range blocks {
		if len(block) < 2 {
			continue
		}
		// Extract text from Tj and TJ operators
		text := extractTextOperators(block[1])
		if text != "" {
			texts = append(texts, text)
		}
	}

	result := strings.Join(texts, "\n")
	result = normalizeWhitespace(result)
	return strings.TrimSpace(result)
}

// extractTextOperators extracts text from PDF text operators (Tj, TJ, ').
func extractTextOperators(block string) string {
	var parts []string

	// Match Tj operator: (text) Tj
	tjRe := regexp.MustCompile(`\(([^)]*)\)\s*Tj`)
	for _, m := range tjRe.FindAllStringSubmatch(block, -1) {
		if len(m) > 1 {
			parts = append(parts, decodePDFString(m[1]))
		}
	}

	// Match TJ operator: [(text) -kern (text)] TJ
	tjArrRe := regexp.MustCompile(`\[([^\]]*)\]\s*TJ`)
	for _, m := range tjArrRe.FindAllStringSubmatch(block, -1) {
		if len(m) > 1 {
			innerRe := regexp.MustCompile(`\(([^)]*)\)`)
			for _, inner := range innerRe.FindAllStringSubmatch(m[1], -1) {
				if len(inner) > 1 {
					parts = append(parts, decodePDFString(inner[1]))
				}
			}
		}
	}

	return strings.Join(parts, "")
}

// decodePDFString decodes basic PDF string escape sequences.
func decodePDFString(s string) string {
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\r", "\r")
	s = strings.ReplaceAll(s, "\\t", "\t")
	s = strings.ReplaceAll(s, "\\(", "(")
	s = strings.ReplaceAll(s, "\\)", ")")
	s = strings.ReplaceAll(s, "\\\\", "\\")
	return s
}
