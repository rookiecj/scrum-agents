package extractor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// YouTubeExtractor extracts transcripts and metadata from YouTube videos.
type YouTubeExtractor struct {
	Client *http.Client
}

// NewYouTubeExtractor creates a new YouTubeExtractor.
func NewYouTubeExtractor() *YouTubeExtractor {
	return &YouTubeExtractor{
		Client: &http.Client{},
	}
}

// VideoMetadata holds YouTube video metadata.
type VideoMetadata struct {
	Title       string `json:"title"`
	Channel     string `json:"channel"`
	Description string `json:"description"`
}

// Extract fetches the transcript from a YouTube video URL.
func (e *YouTubeExtractor) Extract(rawURL string) (*model.ExtractedContent, error) {
	videoID, err := extractVideoID(rawURL)
	if err != nil {
		return nil, fmt.Errorf("extracting video ID: %w", err)
	}

	// Fetch the YouTube page to get metadata and transcript data
	pageURL := "https://www.youtube.com/watch?v=" + videoID
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LinkSummarizer/1.0)")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ko;q=0.8")

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching YouTube page: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	pageHTML := string(body)

	metadata := extractVideoMetadata(pageHTML)

	linkInfo := model.LinkInfo{
		URL:      rawURL,
		LinkType: model.LinkTypeYouTube,
		Title:    metadata.Title,
		Author:   metadata.Channel,
	}

	// Try to get captions URL from the page
	captionsURL, err := extractCaptionsURL(pageHTML)
	if err == nil {
		// Fetch the transcript
		transcript, fetchErr := e.fetchTranscript(captionsURL)
		if fetchErr == nil && transcript != "" {
			return &model.ExtractedContent{
				LinkInfo: linkInfo,
				Content:  transcript,
			}, nil
		}
	}

	// Fallback: build content from video metadata (title + description + channel)
	content := buildMetadataContent(metadata)
	return &model.ExtractedContent{
		LinkInfo: linkInfo,
		Content:  content,
	}, nil
}

// extractVideoID extracts the video ID from various YouTube URL formats.
func extractVideoID(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parsing URL: %w", err)
	}

	host := strings.ToLower(u.Hostname())

	// youtu.be/<id>
	if strings.Contains(host, "youtu.be") {
		id := strings.TrimPrefix(u.Path, "/")
		if id == "" {
			return "", fmt.Errorf("no video ID in short URL")
		}
		return id, nil
	}

	// youtube.com/watch?v=<id>
	if strings.Contains(host, "youtube.com") {
		if v := u.Query().Get("v"); v != "" {
			return v, nil
		}
		// youtube.com/embed/<id> or youtube.com/v/<id>
		parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
		if len(parts) >= 2 && (parts[0] == "embed" || parts[0] == "v") {
			return parts[1], nil
		}
		return "", fmt.Errorf("no video ID found in YouTube URL")
	}

	return "", fmt.Errorf("not a YouTube URL: %s", host)
}

// extractVideoMetadata extracts title, channel, and description from the YouTube page HTML.
func extractVideoMetadata(html string) VideoMetadata {
	meta := VideoMetadata{}

	// Extract title from <meta property="og:title">
	titleRe := regexp.MustCompile(`<meta\s+property="og:title"\s+content="([^"]*)"`)
	if m := titleRe.FindStringSubmatch(html); len(m) > 1 {
		meta.Title = m[1]
	}

	// Extract channel from ownerChannelName
	channelRe := regexp.MustCompile(`"ownerChannelName":"([^"]*)"`)
	if m := channelRe.FindStringSubmatch(html); len(m) > 1 {
		meta.Channel = m[1]
	}

	// Extract full description from shortDescription JSON field
	descRe := regexp.MustCompile(`"shortDescription":"((?:[^"\\]|\\.)*)"`)
	if m := descRe.FindStringSubmatch(html); len(m) > 1 {
		desc := strings.ReplaceAll(m[1], `\n`, "\n")
		desc = strings.ReplaceAll(desc, `\"`, `"`)
		desc = strings.ReplaceAll(desc, `\\`, `\`)
		meta.Description = desc
	}

	return meta
}

// buildMetadataContent creates summarizable content from video metadata
// when transcript extraction is not available.
func buildMetadataContent(meta VideoMetadata) string {
	var parts []string

	if meta.Title != "" {
		parts = append(parts, "Title: "+meta.Title)
	}
	if meta.Channel != "" {
		parts = append(parts, "Channel: "+meta.Channel)
	}
	if meta.Description != "" {
		parts = append(parts, "Description:\n"+meta.Description)
	}

	if len(parts) == 0 {
		return "No content available for this video."
	}

	return strings.Join(parts, "\n\n")
}

// extractCaptionsURL finds the captions/transcript URL from the YouTube page source.
func extractCaptionsURL(html string) (string, error) {
	// Look for captionTracks in the page source
	re := regexp.MustCompile(`"captionTracks":\[(\{[^]]*\})\]`)
	match := re.FindStringSubmatch(html)
	if len(match) < 2 {
		return "", fmt.Errorf("no caption tracks found")
	}

	// Parse the first caption track to get the URL
	trackData := match[1]

	// Extract baseUrl
	urlRe := regexp.MustCompile(`"baseUrl":"([^"]*)"`)
	urlMatch := urlRe.FindStringSubmatch(trackData)
	if len(urlMatch) < 2 {
		return "", fmt.Errorf("no caption URL found")
	}

	captionsURL := strings.ReplaceAll(urlMatch[1], `\u0026`, "&")
	return captionsURL, nil
}

// fetchTranscript downloads and parses the caption XML into plain text.
func (e *YouTubeExtractor) fetchTranscript(captionsURL string) (string, error) {
	// Append fmt=json3 for JSON format, or use XML
	if !strings.Contains(captionsURL, "fmt=") {
		if strings.Contains(captionsURL, "?") {
			captionsURL += "&fmt=json3"
		} else {
			captionsURL += "?fmt=json3"
		}
	}

	resp, err := e.Client.Get(captionsURL)
	if err != nil {
		return "", fmt.Errorf("fetching captions: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading captions: %w", err)
	}

	// Try JSON format first
	transcript, err := parseJSON3Transcript(body)
	if err != nil {
		// Fallback: try XML format
		return parseXMLTranscript(string(body)), nil
	}

	return transcript, nil
}

// parseJSON3Transcript parses YouTube's json3 caption format.
func parseJSON3Transcript(data []byte) (string, error) {
	var result struct {
		Events []struct {
			Segs []struct {
				UTF8 string `json:"utf8"`
			} `json:"segs"`
		} `json:"events"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("parsing json3: %w", err)
	}

	var lines []string
	for _, event := range result.Events {
		for _, seg := range event.Segs {
			text := strings.TrimSpace(seg.UTF8)
			if text != "" && text != "\n" {
				lines = append(lines, text)
			}
		}
	}

	if len(lines) == 0 {
		return "", fmt.Errorf("no transcript content found")
	}

	return strings.Join(lines, " "), nil
}

// parseXMLTranscript extracts text from YouTube's XML caption format.
func parseXMLTranscript(xml string) string {
	// Simple extraction: find all text between <text> tags
	var lines []string
	re := regexp.MustCompile(`<text[^>]*>([^<]*)</text>`)
	matches := re.FindAllStringSubmatch(xml, -1)
	for _, m := range matches {
		if len(m) > 1 {
			text := strings.TrimSpace(m[1])
			// Unescape basic HTML entities
			text = strings.ReplaceAll(text, "&amp;", "&")
			text = strings.ReplaceAll(text, "&lt;", "<")
			text = strings.ReplaceAll(text, "&gt;", ">")
			text = strings.ReplaceAll(text, "&#39;", "'")
			text = strings.ReplaceAll(text, "&quot;", `"`)
			if text != "" {
				lines = append(lines, text)
			}
		}
	}
	return strings.Join(lines, " ")
}
