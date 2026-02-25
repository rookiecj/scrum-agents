package extractor

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTwitterExtractor_Extract(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErr    bool
		errContain string
		wantAuthor string
	}{
		{
			name: "tweet with OG meta tags",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<html><head>
					<meta property="og:title" content="@elonmusk on X">
					<meta property="og:description" content="This is a great tweet about technology!">
					<meta name="twitter:creator" content="@elonmusk">
				</head><body></body></html>`))
			},
			wantAuthor: "@elonmusk",
		},
		{
			name: "tweet with HTML entities",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<html><head>
					<meta property="og:description" content="Hello &amp; welcome &lt;world&gt;">
				</head><body></body></html>`))
			},
		},
		{
			name: "private tweet",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			},
			wantErr:    true,
			errContain: "private or protected",
		},
		{
			name: "unauthorized tweet",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			},
			wantErr:    true,
			errContain: "private or protected",
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
			name: "no extractable content",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("<html><head></head><body></body></html>"))
			},
			wantErr:    true,
			errContain: "could not extract tweet content",
		},
		{
			name: "fallback to body content",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<html><head></head><body>
					<div class="tweet-text">Interesting thread about AI safety</div>
				</body></html>`))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			ext := &TwitterExtractor{Client: server.Client()}
			result, err := ext.Extract(server.URL + "/tweet/123")

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
			if result.LinkInfo.LinkType != "twitter" {
				t.Errorf("LinkType = %q, want %q", result.LinkInfo.LinkType, "twitter")
			}
			if tt.wantAuthor != "" && result.LinkInfo.Author != tt.wantAuthor {
				t.Errorf("Author = %q, want %q", result.LinkInfo.Author, tt.wantAuthor)
			}
		})
	}
}

func TestExtractOGMeta(t *testing.T) {
	html := `<html><head>
		<meta property="og:title" content="My Title">
		<meta property="og:description" content="My Description">
		<meta content="Reversed" property="og:site_name">
		<meta name="twitter:creator" content="@user">
	</head></html>`

	tests := []struct {
		property string
		want     string
	}{
		{"og:title", "My Title"},
		{"og:description", "My Description"},
		{"og:site_name", "Reversed"},
		{"twitter:creator", "@user"},
		{"og:missing", ""},
	}

	for _, tt := range tests {
		t.Run(tt.property, func(t *testing.T) {
			got := extractOGMeta(html, tt.property)
			if got != tt.want {
				t.Errorf("extractOGMeta(%q) = %q, want %q", tt.property, got, tt.want)
			}
		})
	}
}

func TestDecodeHTMLEntities(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello &amp; World", "Hello & World"},
		{"&lt;div&gt;", "<div>"},
		{"&quot;quoted&quot;", "\"quoted\""},
		{"it&#39;s", "it's"},
		{"no entities", "no entities"},
	}

	for _, tt := range tests {
		got := decodeHTMLEntities(tt.input)
		if got != tt.want {
			t.Errorf("decodeHTMLEntities(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
