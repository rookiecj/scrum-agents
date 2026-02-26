package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/rookiecj/scrum-agents/backend/internal/auth"
	"github.com/rookiecj/scrum-agents/backend/internal/handler"
	"github.com/rookiecj/scrum-agents/backend/internal/llm"
	"github.com/rookiecj/scrum-agents/backend/internal/logging"
	"github.com/rookiecj/scrum-agents/backend/internal/summarizer"
)

func main() {
	logging.Init()

	// Load .env file if it exists; fall back to system environment variables otherwise.
	if err := loadEnv(""); err != nil {
		slog.Warn("could not load .env file, using system environment variables",
			slog.String("error", err.Error()),
		)
	}

	// Database & Auth setup
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "scrum-agents.db"
	}

	store, err := auth.NewStore(dbPath)
	if err != nil {
		slog.Error("failed to initialise database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer store.Close()
	slog.Info("database initialised", slog.String("path", dbPath))

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
		slog.Warn("JWT_SECRET not set, using insecure default — set JWT_SECRET for production")
	}
	jwtSvc := auth.NewJWTService(jwtSecret, 24*time.Hour)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","version":"%s"}`, Version)
	})

	// Auth endpoints (public)
	mux.HandleFunc("POST /api/signup", handler.HandleSignup(store))
	mux.HandleFunc("POST /api/login", handler.HandleLogin(store, jwtSvc))

	// Public API endpoints
	mux.HandleFunc("POST /api/detect", handler.HandleDetect())
	mux.HandleFunc("POST /api/extract", handler.HandleExtract())
	mux.HandleFunc("GET /api/providers", handler.HandleProviders())

	// LLM provider registration
	providers := make(map[string]llm.Provider)

	claudeKey := os.Getenv("ANTHROPIC_API_KEY")
	if claudeKey == "" {
		claudeKey = os.Getenv("CLAUDE_API_KEY")
	}
	claudeProvider := llm.NewClaudeProvider(llm.DefaultClaudeConfig(claudeKey))
	providers["claude"] = claudeProvider

	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey != "" {
		providers["openai"] = llm.NewOpenAIProvider(llm.DefaultOpenAIConfig(openaiKey))
		slog.Info("OpenAI provider registered")
	} else {
		slog.Warn("OPENAI_API_KEY not set, OpenAI provider disabled")
	}

	googleKey := os.Getenv("GOOGLE_API_KEY")
	if googleKey != "" {
		providers["gemini"] = llm.NewGeminiProvider(llm.DefaultGeminiConfig(googleKey))
		slog.Info("Gemini provider registered")
	} else {
		slog.Warn("GOOGLE_API_KEY not set, Gemini provider disabled")
	}

	// Use Claude as the default provider for LLM-dependent endpoints
	defaultProvider := claudeProvider

	mux.HandleFunc("POST /api/classify", handler.HandleClassify(defaultProvider, providers))

	registry, err := summarizer.LoadTemplates("prompts")
	if err != nil {
		slog.Warn("could not load prompt templates, summarize endpoint disabled",
			slog.String("error", err.Error()),
		)
	} else {
		sum := summarizer.NewSummarizer(registry, 0.6)
		mux.HandleFunc("POST /api/summarize", handler.HandleSummarize(sum, defaultProvider, providers))
		slog.Info("prompt templates loaded",
			slog.Int("template_count", len(registry.Categories())),
		)
	}

	if claudeKey == "" {
		slog.Warn("ANTHROPIC_API_KEY not set, LLM endpoints will return errors")
	}

	addr := ":8080"
	slog.Info("starting server", slog.String("addr", addr), slog.String("version", Version))
	if err := http.ListenAndServe(addr, logging.Middleware(mux)); err != nil {
		slog.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
