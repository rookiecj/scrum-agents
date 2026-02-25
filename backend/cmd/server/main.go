package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/rookiecj/scrum-agents/backend/internal/handler"
	"github.com/rookiecj/scrum-agents/backend/internal/llm"
	"github.com/rookiecj/scrum-agents/backend/internal/logging"
	"github.com/rookiecj/scrum-agents/backend/internal/summarizer"
)

func main() {
	logging.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","version":"%s"}`, Version)
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
		slog.Warn("could not load prompt templates, summarize endpoint disabled",
			slog.String("error", err.Error()),
		)
	} else {
		sum := summarizer.NewSummarizer(registry, 0.6)
		mux.HandleFunc("POST /api/summarize", handler.HandleSummarize(sum, provider))
		slog.Info("prompt templates loaded",
			slog.Int("template_count", len(registry.Categories())),
		)
	}

	if apiKey == "" {
		slog.Warn("ANTHROPIC_API_KEY not set, LLM endpoints will return errors")
	}

	addr := ":8080"
	slog.Info("starting server", slog.String("addr", addr), slog.String("version", Version))
	if err := http.ListenAndServe(addr, logging.Middleware(mux)); err != nil {
		slog.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
