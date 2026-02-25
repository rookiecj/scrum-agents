package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rookiecj/scrum-agents/backend/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})
	mux.HandleFunc("POST /api/detect", handler.HandleDetect())
	mux.HandleFunc("POST /api/extract", handler.HandleExtract())

	addr := ":8080"
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
