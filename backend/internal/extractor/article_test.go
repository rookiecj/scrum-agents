package extractor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

func TestArticleExtractor_Extract(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		statusCode  int
		wantTitle   string
		wantContent string
		wantErr     bool
	}{
		{
			name: "basic article",
			html: `<html><head><title>Test Article</title></head>
<body><article><p>Hello world. This is a test article.</p></article></body></html>`,
			statusCode:  200,
			wantTitle:   "Test Article",
			wantContent: "Hello world. This is a test article.",
		},
		{
			name: "article with nav and footer stripped",
			html: `<html><head><title>Clean Article</title></head>
<body><nav>Menu items</nav><main><p>Main content here.</p></main><footer>Footer info</footer></body></html>`,
			statusCode:  200,
			wantTitle:   "Clean Article",
			wantContent: "Main content here.",
		},
		{
			name: "article with script stripped",
			html: `<html><head><title>No Scripts</title></head>
<body><p>Real content.</p><script>alert('bad')</script></body></html>`,
			statusCode:  200,
			wantTitle:   "No Scripts",
			wantContent: "Real content.",
		},
		{
			name:       "server error",
			html:       "",
			statusCode: 500,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.html))
			}))
			defer server.Close()

			ext := &ArticleExtractor{Client: server.Client()}
			result, err := ext.Extract(server.URL)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.LinkInfo.Title != tt.wantTitle {
				t.Errorf("title = %q, want %q", result.LinkInfo.Title, tt.wantTitle)
			}
			if result.Content != tt.wantContent {
				t.Errorf("content = %q, want %q", result.Content, tt.wantContent)
			}
			if result.LinkInfo.LinkType != model.LinkTypeArticle {
				t.Errorf("link type = %q, want %q", result.LinkInfo.LinkType, model.LinkTypeArticle)
			}
		})
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{name: "simple title", html: "<title>Hello</title>", want: "Hello"},
		{name: "title with attrs", html: `<title lang="en">Greetings</title>`, want: "Greetings"},
		{name: "no title", html: "<html><body>no title</body></html>", want: ""},
		{name: "empty title", html: "<title></title>", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTitle(tt.html)
			if got != tt.want {
				t.Errorf("extractTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStripTags(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{name: "simple", html: "<p>hello</p>", want: "hello"},
		{name: "nested", html: "<div><p>nested</p></div>", want: "nested"},
		{name: "no tags", html: "plain text", want: "plain text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripTags(tt.html)
			if got != tt.want {
				t.Errorf("stripTags() = %q, want %q", got, tt.want)
			}
		})
	}
}
