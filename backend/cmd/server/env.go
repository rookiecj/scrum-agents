package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// loadEnv loads environment variables from the given .env file path.
// If no path is provided, it defaults to ".env" in the current directory.
// It returns nil if the file was loaded successfully, or if the file does not exist.
// Non-file-not-found errors are returned as warnings (logged but not fatal).
func loadEnv(path string) error {
	var err error
	if path == "" {
		err = godotenv.Load()
	} else {
		err = godotenv.Load(path)
	}

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Info(".env file not found, using system environment variables")
			return nil
		}
		return fmt.Errorf("loading .env file: %w", err)
	}

	slog.Info(".env file loaded successfully")
	return nil
}
