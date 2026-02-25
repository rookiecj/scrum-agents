package urldetect

import (
	"testing"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected model.LinkType
		wantErr  bool
	}{
		// YouTube
		{name: "youtube watch", url: "https://www.youtube.com/watch?v=abc123", expected: model.LinkTypeYouTube},
		{name: "youtube short url", url: "https://youtu.be/abc123", expected: model.LinkTypeYouTube},
		{name: "youtube embed", url: "https://www.youtube.com/embed/abc123", expected: model.LinkTypeYouTube},

		// Twitter/X
		{name: "twitter", url: "https://twitter.com/user/status/123", expected: model.LinkTypeTwitter},
		{name: "x.com", url: "https://x.com/user/status/123", expected: model.LinkTypeTwitter},

		// PDF
		{name: "pdf link", url: "https://example.com/paper.pdf", expected: model.LinkTypePDF},
		{name: "pdf with path", url: "https://arxiv.org/pdf/2301.00001.pdf", expected: model.LinkTypePDF},

		// Newsletter
		{name: "substack", url: "https://newsletter.substack.com/p/some-post", expected: model.LinkTypeNewsletter},
		{name: "medium", url: "https://medium.com/@user/some-article-123", expected: model.LinkTypeNewsletter},

		// Article (default)
		{name: "generic article", url: "https://example.com/blog/some-post", expected: model.LinkTypeArticle},
		{name: "tech blog", url: "https://blog.golang.org/go1.22", expected: model.LinkTypeArticle},
		{name: "news site", url: "https://news.ycombinator.com/item?id=123", expected: model.LinkTypeArticle},

		// Error
		{name: "invalid url", url: "://invalid", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Detect(tt.url)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Detect(%q) expected error, got nil", tt.url)
				}
				return
			}
			if err != nil {
				t.Fatalf("Detect(%q) unexpected error: %v", tt.url, err)
			}
			if got != tt.expected {
				t.Errorf("Detect(%q) = %q, want %q", tt.url, got, tt.expected)
			}
		})
	}
}
