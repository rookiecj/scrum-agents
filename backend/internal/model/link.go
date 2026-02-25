package model

// LinkType represents the detected type of a URL.
type LinkType string

const (
	LinkTypeArticle    LinkType = "article"
	LinkTypeYouTube    LinkType = "youtube"
	LinkTypePDF        LinkType = "pdf"
	LinkTypeTwitter    LinkType = "twitter"
	LinkTypeNewsletter LinkType = "newsletter"
	LinkTypeUnknown    LinkType = "unknown"
)

// LinkInfo holds the result of URL type detection and metadata extraction.
type LinkInfo struct {
	URL      string   `json:"url"`
	LinkType LinkType `json:"link_type"`
	Title    string   `json:"title,omitempty"`
	Author   string   `json:"author,omitempty"`
	Date     string   `json:"date,omitempty"`
}

// ExtractedContent holds the content extracted from a URL.
type ExtractedContent struct {
	LinkInfo LinkInfo `json:"link_info"`
	Content  string  `json:"content"`
}
