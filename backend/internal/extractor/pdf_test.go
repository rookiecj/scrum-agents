package extractor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPDFExtractor_Extract(t *testing.T) {
	// Create a minimal PDF with extractable text
	simplePDF := buildSimplePDF("Hello World from PDF")

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErr    bool
		errContain string
		wantTitle  string
	}{
		{
			name: "basic PDF with text",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/pdf")
				w.Write(simplePDF)
			},
		},
		{
			name: "PDF too large",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/pdf")
				w.Header().Set("Content-Length", fmt.Sprintf("%d", maxPDFSize+1))
				w.Write([]byte("%PDF-1.4"))
			},
			wantErr:    true,
			errContain: "exceeds maximum size",
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
			name: "no extractable text",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/pdf")
				w.Write([]byte("%PDF-1.4\n%%EOF"))
			},
			wantErr:    true,
			errContain: "could not extract text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			ext := &PDFExtractor{Client: server.Client()}
			result, err := ext.Extract(server.URL + "/test.pdf")

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
			if result.LinkInfo.LinkType != "pdf" {
				t.Errorf("LinkType = %q, want %q", result.LinkInfo.LinkType, "pdf")
			}
		})
	}
}

func TestExtractPDFText(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantLen bool // whether we expect some text
	}{
		{
			name:    "simple Tj text",
			data:    []byte("BT (Hello World) Tj ET"),
			wantLen: true,
		},
		{
			name:    "TJ array text",
			data:    []byte("BT [(Hello) -100 ( World)] TJ ET"),
			wantLen: true,
		},
		{
			name:    "no text blocks",
			data:    []byte("%PDF-1.4 some binary data"),
			wantLen: false,
		},
		{
			name:    "escaped characters",
			data:    []byte(`BT (Hello\nWorld) Tj ET`),
			wantLen: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := extractPDFText(tt.data)
			if tt.wantLen && text == "" {
				t.Error("expected non-empty text")
			}
			if !tt.wantLen && text != "" {
				t.Errorf("expected empty text, got %q", text)
			}
		})
	}
}

func TestExtractPDFTitle(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{
			name: "has title",
			data: []byte("/Title (My Document)"),
			want: "My Document",
		},
		{
			name: "no title",
			data: []byte("%PDF-1.4"),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractPDFTitle(tt.data)
			if got != tt.want {
				t.Errorf("extractPDFTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecodePDFString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`Hello\nWorld`, "Hello\nWorld"},
		{`Open \(paren\)`, "Open (paren)"},
		{`Back\\slash`, "Back\\slash"},
		{"No escapes", "No escapes"},
	}

	for _, tt := range tests {
		got := decodePDFString(tt.input)
		if got != tt.want {
			t.Errorf("decodePDFString(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// buildSimplePDF creates a minimal valid PDF with extractable text.
func buildSimplePDF(text string) []byte {
	return []byte(fmt.Sprintf(`%%PDF-1.4
1 0 obj <</Type /Catalog /Pages 2 0 R>> endobj
2 0 obj <</Type /Pages /Kids [3 0 R] /Count 1>> endobj
3 0 obj <</Type /Page /Parent 2 0 R /MediaBox [0 0 612 792]
/Contents 4 0 R /Resources <</Font <</F1 5 0 R>>>>>> endobj
4 0 obj <</Length 44>>
stream
BT /F1 12 Tf 100 700 Td (%s) Tj ET
endstream endobj
5 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica>> endobj
/Title (Test PDF Document)
xref
0 6
trailer <</Size 6 /Root 1 0 R /Info <</Title (Test PDF Document)>>>>
startxref
0
%%%%EOF`, text))
}
