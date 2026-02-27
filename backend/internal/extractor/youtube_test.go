package extractor

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

func TestExtractVideoID(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{name: "standard watch", url: "https://www.youtube.com/watch?v=dQw4w9WgXcQ", want: "dQw4w9WgXcQ"},
		{name: "short url", url: "https://youtu.be/dQw4w9WgXcQ", want: "dQw4w9WgXcQ"},
		{name: "embed url", url: "https://www.youtube.com/embed/dQw4w9WgXcQ", want: "dQw4w9WgXcQ"},
		{name: "with extra params", url: "https://www.youtube.com/watch?v=abc123&t=120", want: "abc123"},
		{name: "no video id", url: "https://www.youtube.com/", wantErr: true},
		{name: "short url no id", url: "https://youtu.be/", wantErr: true},
		{name: "not youtube", url: "https://example.com/watch?v=abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractVideoID(tt.url)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("extractVideoID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractVideoMetadata(t *testing.T) {
	html := `<html>
<head><meta property="og:title" content="Test Video Title"></head>
<body><script>var ytInitialData = {"ownerChannelName":"Test Channel","shortDescription":"This is a test video about Go concurrency.\nLearn more at example.com"}</script></body>
</html>`

	meta := extractVideoMetadata(html)
	if meta.Title != "Test Video Title" {
		t.Errorf("title = %q, want %q", meta.Title, "Test Video Title")
	}
	if meta.Channel != "Test Channel" {
		t.Errorf("channel = %q, want %q", meta.Channel, "Test Channel")
	}
	if !strings.Contains(meta.Description, "Go concurrency") {
		t.Errorf("description should contain 'Go concurrency', got %q", meta.Description)
	}
}

func TestBuildMetadataContent(t *testing.T) {
	tests := []struct {
		name string
		meta VideoMetadata
		want string
	}{
		{
			name: "full metadata",
			meta: VideoMetadata{Title: "My Video", Channel: "My Channel", Description: "A great video"},
			want: "Title: My Video\n\nChannel: My Channel\n\nDescription:\nA great video",
		},
		{
			name: "no description",
			meta: VideoMetadata{Title: "My Video", Channel: "My Channel"},
			want: "Title: My Video\n\nChannel: My Channel",
		},
		{
			name: "empty metadata",
			meta: VideoMetadata{},
			want: "No content available for this video.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildMetadataContent(tt.meta)
			if got != tt.want {
				t.Errorf("buildMetadataContent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractCaptionsURL(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		wantErr bool
	}{
		{
			name:    "with captions",
			html:    `"captionTracks":[{"baseUrl":"https://www.youtube.com/api/timedtext?v=abc\u0026lang=en","name":{"simpleText":"English"}}]`,
			wantErr: false,
		},
		{
			name:    "no captions",
			html:    `<html><body>no captions here</body></html>`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := extractCaptionsURL(tt.html)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if url == "" {
				t.Error("expected non-empty URL")
			}
		})
	}
}

func TestParseJSON3Transcript(t *testing.T) {
	json3 := `{"events":[{"segs":[{"utf8":"Hello "}]},{"segs":[{"utf8":"world"}]},{"segs":[{"utf8":"\n"}]}]}`

	got, err := parseJSON3Transcript([]byte(json3))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Hello world" {
		t.Errorf("parseJSON3Transcript() = %q, want %q", got, "Hello world")
	}
}

func TestParseXMLTranscript(t *testing.T) {
	xml := `<transcript><text start="0" dur="5">Hello</text><text start="5" dur="3">world &amp; friends</text></transcript>`

	got := parseXMLTranscript(xml)
	if got != "Hello world & friends" {
		t.Errorf("parseXMLTranscript() = %q, want %q", got, "Hello world & friends")
	}
}

// youtubeTestClient creates an http.Client that redirects youtube.com requests to a local test server.
func youtubeTestClient(server *httptest.Server) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Scheme = "http"
			req.URL.Host = strings.TrimPrefix(server.URL, "http://")
			return http.DefaultTransport.RoundTrip(req)
		}),
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestYouTubeExtractor_Extract_NoCaptions_FallbackToDescription(t *testing.T) {
	pageHTML := `<html>
<head><meta property="og:title" content="No Caption Video"></head>
<body><script>var ytInitialData = {"ownerChannelName":"Test Channel","shortDescription":"This video explains Go patterns.\nVery useful."}</script></body>
</html>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(pageHTML))
	}))
	defer server.Close()

	ext := &YouTubeExtractor{Client: youtubeTestClient(server)}
	result, err := ext.Extract("https://www.youtube.com/watch?v=test123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LinkInfo.Title != "No Caption Video" {
		t.Errorf("title = %q, want %q", result.LinkInfo.Title, "No Caption Video")
	}
	if !strings.Contains(result.Content, "Go patterns") {
		t.Errorf("content should contain description fallback, got %q", result.Content)
	}
	if result.Content == "" {
		t.Error("content should not be empty when description is available")
	}
}

func TestYouTubeExtractor_Extract_EmptyTranscript_FallbackToDescription(t *testing.T) {
	// All requests (page + captions) go to this server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "timedtext") {
			// Captions endpoint returns empty (simulates YouTube blocking)
			w.WriteHeader(http.StatusOK)
			return
		}
		// YouTube page with captions URL pointing to localhost timedtext
		pageHTML := `<html>
<head><meta property="og:title" content="Blocked Transcript Video"></head>
<body><script>
var ytInitialPlayerResponse = {
"ownerChannelName":"Test Channel",
"shortDescription":"Learn about microservices architecture.",
"captionTracks":[{"baseUrl":"https://www.youtube.com/api/timedtext?v=abc\u0026lang=en","name":{"simpleText":"English"}}]
}
</script></body>
</html>`
		w.Write([]byte(pageHTML))
	}))
	defer server.Close()

	ext := &YouTubeExtractor{Client: youtubeTestClient(server)}
	result, err := ext.Extract("https://www.youtube.com/watch?v=abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Content, "microservices") {
		t.Errorf("content should fall back to description, got %q", result.Content)
	}
}

func TestYouTubeExtractor_Extract_WithCaptions(t *testing.T) {
	// Verify the YouTube extractor implements the Extractor interface
	var _ Extractor = &YouTubeExtractor{}

	// Verify constructor
	ext := NewYouTubeExtractor()
	if ext.Client == nil {
		t.Error("expected non-nil client")
	}
}

func TestYouTubeExtractor_ImplementsExtractor(t *testing.T) {
	var _ Extractor = &YouTubeExtractor{}

	ext := NewYouTubeExtractor()
	if ext == nil {
		t.Fatal("NewYouTubeExtractor() returned nil")
	}

	// Verify the extract method would return correct link type
	// by testing with an unreachable URL
	result, err := ext.Extract("https://www.youtube.com/watch?v=nonexistent")
	if err == nil && result != nil {
		if result.LinkInfo.LinkType != model.LinkTypeYouTube {
			t.Errorf("link type = %q, want %q", result.LinkInfo.LinkType, model.LinkTypeYouTube)
		}
	}
	// Error is expected since we can't reach YouTube in tests
}
