package extractor

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// NewsletterExtractor extracts content from newsletter platforms (Substack, Medium, etc.).
type NewsletterExtractor struct {
	Client *http.Client
}

// NewNewsletterExtractor creates a new NewsletterExtractor.
func NewNewsletterExtractor() *NewsletterExtractor {
	return &NewsletterExtractor{
		Client: &http.Client{},
	}
}

// Extract fetches a newsletter page and extracts the article body.
func (e *NewsletterExtractor) Extract(rawURL string) (*model.ExtractedContent, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LinkSummarizer/1.0)")

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching newsletter %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for newsletter %s", resp.StatusCode, rawURL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	html := string(body)
	title := extractTitle(html)
	author := extractNewsletterAuthor(html)
	content := extractNewsletterContent(html)

	partialExtraction := isPaywalled(html)

	if content == "" {
		// Fallback to generic article extraction
		content = extractMainContent(html)
	}

	if content == "" {
		return nil, fmt.Errorf("could not extract newsletter content")
	}

	if partialExtraction {
		content = content + "\n\n---\n[Note: This content may be partially extracted due to paywall restrictions]"
	}

	return &model.ExtractedContent{
		LinkInfo: model.LinkInfo{
			URL:      rawURL,
			LinkType: model.LinkTypeNewsletter,
			Title:    title,
			Author:   author,
		},
		Content: content,
	}, nil
}

// extractNewsletterContent extracts the main article content from newsletter platforms.
func extractNewsletterContent(html string) string {
	// Try platform-specific content containers first
	containers := []struct {
		open  string
		close string
	}{
		// Substack
		{`<div class="body markup"`, `</div>`},
		{`<div class="post-content"`, `</div>`},
		// Medium
		{`<article`, `</article>`},
		{`<div class="section-content"`, `</div>`},
	}

	for _, c := range containers {
		content := extractBlock(html, c.open, c.close)
		if content != "" {
			content = stripTags(content)
			content = normalizeWhitespace(content)
			if len(content) > 100 { // Sanity check: meaningful content
				return strings.TrimSpace(content)
			}
		}
	}

	return ""
}

// extractBlock extracts content between an opening tag pattern and its closing tag.
func extractBlock(html, openPattern, closeTag string) string {
	lower := strings.ToLower(html)
	openLower := strings.ToLower(openPattern)

	start := strings.Index(lower, openLower)
	if start == -1 {
		return ""
	}

	// Find the end of the opening tag
	tagEnd := strings.Index(html[start:], ">")
	if tagEnd == -1 {
		return ""
	}
	contentStart := start + tagEnd + 1

	// Find the matching close tag (simplified - finds the next occurrence)
	closeLower := strings.ToLower(closeTag)
	end := strings.Index(lower[contentStart:], closeLower)
	if end == -1 {
		return ""
	}

	return html[contentStart : contentStart+end]
}

// extractNewsletterAuthor extracts the author name from newsletter meta tags.
func extractNewsletterAuthor(html string) string {
	author := extractOGMeta(html, "author")
	if author != "" {
		return author
	}
	author = extractOGMeta(html, "article:author")
	if author != "" {
		return author
	}
	return extractOGMeta(html, "og:site_name")
}

// isPaywalled checks if the newsletter content appears to be behind a paywall.
func isPaywalled(html string) bool {
	lower := strings.ToLower(html)
	indicators := []string{
		"paywall",
		"subscribe to continue",
		"members-only",
		"premium content",
		"upgrade to read",
		"this post is for paid subscribers",
		"this post is for paying subscribers",
	}
	for _, indicator := range indicators {
		if strings.Contains(lower, indicator) {
			return true
		}
	}
	return false
}
