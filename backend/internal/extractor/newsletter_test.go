package extractor

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewsletterExtractor_Extract(t *testing.T) {
	substackHTML := `<html><head>
		<title>My Newsletter Post</title>
		<meta property="author" content="John Doe">
	</head><body>
		<div class="body markup">` + strings.Repeat("This is the newsletter content. ", 10) + `</div>
	</body></html>`

	mediumHTML := `<html><head>
		<title>Medium Article</title>
		<meta property="article:author" content="Jane Smith">
	</head><body>
		<article>` + strings.Repeat("This is a medium article with lots of content. ", 10) + `</article>
	</body></html>`

	paywallHTML := `<html><head><title>Premium Post</title></head><body>
		<div class="body markup">` + strings.Repeat("Some content visible before paywall. ", 10) + `</div>
		<div class="paywall">Subscribe to continue reading</div>
	</body></html>`

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		errContain  string
		wantPaywall bool
		wantAuthor  string
	}{
		{
			name: "substack newsletter",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(substackHTML))
			},
			wantAuthor: "John Doe",
		},
		{
			name: "medium article",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(mediumHTML))
			},
			wantAuthor: "Jane Smith",
		},
		{
			name: "paywalled newsletter",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(paywallHTML))
			},
			wantPaywall: true,
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr:    true,
			errContain: "unexpected status",
		},
		{
			name: "empty page",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("<html><head></head><body></body></html>"))
			},
			wantErr:    true,
			errContain: "could not extract",
		},
		{
			name: "fallback to generic extraction",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<html><head><title>Unknown Newsletter</title></head>
					<body><p>` + strings.Repeat("Fallback content from unknown platform. ", 10) + `</p></body></html>`))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			ext := &NewsletterExtractor{Client: server.Client()}
			result, err := ext.Extract(server.URL + "/newsletter/post")

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errContain != "" && !strings.Contains(err.Error(), tt.errContain) {
					t.Errorf("error = %q, want containing %q", err.Error(), tt.errContain)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Content == "" {
				t.Error("expected non-empty content")
			}
			if result.LinkInfo.LinkType != "newsletter" {
				t.Errorf("LinkType = %q, want %q", result.LinkInfo.LinkType, "newsletter")
			}
			if tt.wantPaywall && !strings.Contains(result.Content, "partially extracted") {
				t.Error("expected paywall indicator in content")
			}
			if tt.wantAuthor != "" && result.LinkInfo.Author != tt.wantAuthor {
				t.Errorf("Author = %q, want %q", result.LinkInfo.Author, tt.wantAuthor)
			}
		})
	}
}

func TestExtractNewsletterContent(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		wantLen bool
	}{
		{
			name:    "substack body markup",
			html:    `<div class="body markup">` + strings.Repeat("Newsletter content here. ", 10) + `</div>`,
			wantLen: true,
		},
		{
			name:    "medium article tag",
			html:    `<article>` + strings.Repeat("Medium article text here. ", 10) + `</article>`,
			wantLen: true,
		},
		{
			name:    "no matching container",
			html:    `<div>Some random content</div>`,
			wantLen: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractNewsletterContent(tt.html)
			if tt.wantLen && got == "" {
				t.Error("expected non-empty content")
			}
			if !tt.wantLen && got != "" {
				t.Errorf("expected empty content, got %q", got)
			}
		})
	}
}

func TestIsPaywalled(t *testing.T) {
	tests := []struct {
		name string
		html string
		want bool
	}{
		{"has paywall", `<div class="paywall">Subscribe</div>`, true},
		{"members only", `<span>This is members-only content</span>`, true},
		{"premium content", `<p>Premium Content</p>`, true},
		{"subscribe to continue", `<div>Subscribe to continue reading</div>`, true},
		{"no paywall", `<div>Free content for all</div>`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPaywalled(tt.html)
			if got != tt.want {
				t.Errorf("isPaywalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractBlock(t *testing.T) {
	tests := []struct {
		name  string
		html  string
		open  string
		close string
		want  string
	}{
		{
			name:  "simple block",
			html:  `<div class="content">Hello World</div>`,
			open:  `<div class="content"`,
			close: `</div>`,
			want:  "Hello World",
		},
		{
			name:  "no match",
			html:  `<div>Hello</div>`,
			open:  `<article`,
			close: `</article>`,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractBlock(tt.html, tt.open, tt.close)
			if got != tt.want {
				t.Errorf("extractBlock() = %q, want %q", got, tt.want)
			}
		})
	}
}
