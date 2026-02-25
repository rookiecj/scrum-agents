package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rookiecj/scrum-agents/backend/internal/handler"
	"github.com/rookiecj/scrum-agents/backend/internal/llm"
	"github.com/rookiecj/scrum-agents/backend/internal/summarizer"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})
	mux.HandleFunc("POST /api/detect", handler.HandleDetect())
	mux.HandleFunc("POST /api/extract", handler.HandleExtract())

	// LLM-dependent endpoints
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("CLAUDE_API_KEY")
	}
	provider := llm.NewClaudeProvider(llm.DefaultClaudeConfig(apiKey))

	mux.HandleFunc("POST /api/classify", handler.HandleClassify(provider))

	registry, err := summarizer.LoadTemplates("prompts")
	if err != nil {
		log.Printf("Warning: could not load prompt templates: %v (summarize endpoint disabled)", err)
	} else {
		sum := summarizer.NewSummarizer(registry, 0.6)
		mux.HandleFunc("POST /api/summarize", handler.HandleSummarize(sum, provider))
		log.Printf("Loaded %d prompt templates + generic fallback", len(registry.Categories()))
	}

	if apiKey == "" {
		log.Printf("Warning: ANTHROPIC_API_KEY not set, LLM endpoints will return errors")
	}

	addr := ":8080"
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
