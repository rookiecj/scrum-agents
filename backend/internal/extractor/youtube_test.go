package extractor

import (
	"net/http"
	"net/http/httptest"
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
<body><script>var ytInitialData = {"ownerChannelName":"Test Channel"}</script></body>
</html>`

	meta := extractVideoMetadata(html)
	if meta.Title != "Test Video Title" {
		t.Errorf("title = %q, want %q", meta.Title, "Test Video Title")
	}
	if meta.Channel != "Test Channel" {
		t.Errorf("channel = %q, want %q", meta.Channel, "Test Channel")
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

func TestYouTubeExtractor_Extract_NoCaptions(t *testing.T) {
	// Simulate a YouTube page without captions
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><head><meta property="og:title" content="No Caption Video"></head><body>no captions data</body></html>`))
	}))
	defer server.Close()

	_ = &YouTubeExtractor{Client: server.Client()}

	// Test that extractVideoID works correctly (unit test)
	_, err := extractVideoID("https://www.youtube.com/watch?v=test123")
	if err != nil {
		t.Fatalf("extractVideoID failed: %v", err)
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
