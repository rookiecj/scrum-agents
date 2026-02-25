package extractor

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// ArticleExtractor extracts content from web articles using HTML parsing.
type ArticleExtractor struct {
	Client *http.Client
}

// NewArticleExtractor creates a new ArticleExtractor with a default HTTP client.
func NewArticleExtractor() *ArticleExtractor {
	return &ArticleExtractor{
		Client: &http.Client{},
	}
}

// Extract fetches and extracts the main content from a web article URL.
func (e *ArticleExtractor) Extract(rawURL string) (*model.ExtractedContent, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LinkSummarizer/1.0)")

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching URL %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for URL %s", resp.StatusCode, rawURL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	html := string(body)
	title := extractTitle(html)
	content := extractMainContent(html)

	return &model.ExtractedContent{
		LinkInfo: model.LinkInfo{
			URL:      rawURL,
			LinkType: model.LinkTypeArticle,
			Title:    title,
		},
		Content: content,
	}, nil
}

// extractTitle extracts the <title> tag content from HTML.
func extractTitle(html string) string {
	start := strings.Index(html, "<title")
	if start == -1 {
		return ""
	}
	// Skip past the closing > of the opening tag
	start = strings.Index(html[start:], ">")
	if start == -1 {
		return ""
	}
	// Adjust start to be absolute position after >
	titleStart := strings.Index(html, "<title")
	start = titleStart + start + 1

	end := strings.Index(html[start:], "</title>")
	if end == -1 {
		return ""
	}
	return strings.TrimSpace(html[start : start+end])
}

// extractMainContent extracts text from the HTML body, removing tags and scripts.
func extractMainContent(html string) string {
	// Remove script and style blocks
	content := removeBlocks(html, "script")
	content = removeBlocks(content, "style")
	content = removeBlocks(content, "nav")
	content = removeBlocks(content, "header")
	content = removeBlocks(content, "footer")

	// Extract body content
	bodyStart := strings.Index(strings.ToLower(content), "<body")
	if bodyStart != -1 {
		bodyStart = strings.Index(content[bodyStart:], ">")
		if bodyStart != -1 {
			absBodyStart := strings.Index(strings.ToLower(content), "<body")
			bodyStart = absBodyStart + bodyStart + 1
		}
		bodyEnd := strings.Index(strings.ToLower(content[bodyStart:]), "</body>")
		if bodyEnd != -1 {
			content = content[bodyStart : bodyStart+bodyEnd]
		} else {
			content = content[bodyStart:]
		}
	}

	// Strip remaining HTML tags
	content = stripTags(content)

	// Normalize whitespace
	content = normalizeWhitespace(content)

	return strings.TrimSpace(content)
}

// removeBlocks removes all occurrences of a block-level tag and its contents.
func removeBlocks(html, tag string) string {
	result := html
	lowerResult := strings.ToLower(result)
	openTag := "<" + tag
	closeTag := "</" + tag + ">"

	for {
		start := strings.Index(lowerResult, openTag)
		if start == -1 {
			break
		}
		end := strings.Index(lowerResult[start:], closeTag)
		if end == -1 {
			break
		}
		end = start + end + len(closeTag)
		result = result[:start] + result[end:]
		lowerResult = strings.ToLower(result)
	}
	return result
}

// stripTags removes all HTML tags from a string.
func stripTags(html string) string {
	var result strings.Builder
	inTag := false
	for _, r := range html {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			result.WriteRune(r)
		}
	}
	return result.String()
}

// normalizeWhitespace collapses multiple whitespace characters into single spaces
// and multiple newlines into double newlines.
func normalizeWhitespace(s string) string {
	// Split into lines and process
	lines := strings.Split(s, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.Join(strings.Fields(line), " ")
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, "\n")
}
