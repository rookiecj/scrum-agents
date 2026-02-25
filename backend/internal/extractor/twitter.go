package extractor

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// TwitterExtractor extracts tweet content from Twitter/X URLs.
type TwitterExtractor struct {
	Client *http.Client
}

// NewTwitterExtractor creates a new TwitterExtractor.
func NewTwitterExtractor() *TwitterExtractor {
	return &TwitterExtractor{
		Client: &http.Client{},
	}
}

// Extract fetches a tweet page and extracts the tweet text and metadata.
func (e *TwitterExtractor) Extract(rawURL string) (*model.ExtractedContent, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LinkSummarizer/1.0)")

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching tweet %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("tweet is private or protected (status: %d)", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for tweet %s", resp.StatusCode, rawURL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	html := string(body)
	title := extractOGMeta(html, "og:title")
	description := extractOGMeta(html, "og:description")
	author := extractTweetAuthor(html)

	content := description
	if content == "" {
		// Fallback: try to extract from page content
		content = extractMainContent(html)
	}

	if content == "" {
		return nil, fmt.Errorf("could not extract tweet content (tweet may be private or protected)")
	}

	return &model.ExtractedContent{
		LinkInfo: model.LinkInfo{
			URL:      rawURL,
			LinkType: model.LinkTypeTwitter,
			Title:    title,
			Author:   author,
		},
		Content: content,
	}, nil
}

// extractOGMeta extracts Open Graph meta tag content.
func extractOGMeta(html, property string) string {
	// Match <meta property="og:..." content="...">
	patterns := []string{
		fmt.Sprintf(`<meta\s+property="%s"\s+content="([^"]*)"`, regexp.QuoteMeta(property)),
		fmt.Sprintf(`<meta\s+content="([^"]*)"\s+property="%s"`, regexp.QuoteMeta(property)),
		fmt.Sprintf(`<meta\s+name="%s"\s+content="([^"]*)"`, regexp.QuoteMeta(property)),
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 {
			return decodeHTMLEntities(matches[1])
		}
	}
	return ""
}

// extractTweetAuthor extracts the tweet author from meta tags.
func extractTweetAuthor(html string) string {
	// Try twitter:creator or og:site_name
	author := extractOGMeta(html, "twitter:creator")
	if author != "" {
		return author
	}
	return extractOGMeta(html, "og:site_name")
}

// decodeHTMLEntities decodes common HTML entities.
func decodeHTMLEntities(s string) string {
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&apos;", "'")
	return s
}
