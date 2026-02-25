package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnv_FileExists(t *testing.T) {
	// Create a temporary .env file
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := "TEST_LOAD_ENV_VAR=hello_from_dotenv\n"
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	// Clear the variable first
	os.Unsetenv("TEST_LOAD_ENV_VAR")

	err := loadEnv(envPath)
	if err != nil {
		t.Fatalf("loadEnv() returned error: %v", err)
	}

	got := os.Getenv("TEST_LOAD_ENV_VAR")
	if got != "hello_from_dotenv" {
		t.Errorf("TEST_LOAD_ENV_VAR = %q, want %q", got, "hello_from_dotenv")
	}

	// Cleanup
	os.Unsetenv("TEST_LOAD_ENV_VAR")
}

func TestLoadEnv_FileNotExists(t *testing.T) {
	// Point to a non-existent file
	err := loadEnv("/nonexistent/path/.env")
	if err != nil {
		t.Errorf("loadEnv() should not return error for missing file, got: %v", err)
	}
}

func TestLoadEnv_DoesNotOverrideExisting(t *testing.T) {
	// Set the variable before loading
	os.Setenv("TEST_NO_OVERRIDE_VAR", "original_value")

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := "TEST_NO_OVERRIDE_VAR=overridden_value\n"
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	err := loadEnv(envPath)
	if err != nil {
		t.Fatalf("loadEnv() returned error: %v", err)
	}

	// godotenv should NOT override existing env vars
	got := os.Getenv("TEST_NO_OVERRIDE_VAR")
	if got != "original_value" {
		t.Errorf("TEST_NO_OVERRIDE_VAR = %q, want %q (should not be overridden)", got, "original_value")
	}

	// Cleanup
	os.Unsetenv("TEST_NO_OVERRIDE_VAR")
}

func TestLoadEnv_MultipleVars(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := `TEST_MULTI_A=value_a
TEST_MULTI_B=value_b
TEST_MULTI_C=value_c
`
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	// Clear variables first
	os.Unsetenv("TEST_MULTI_A")
	os.Unsetenv("TEST_MULTI_B")
	os.Unsetenv("TEST_MULTI_C")

	err := loadEnv(envPath)
	if err != nil {
		t.Fatalf("loadEnv() returned error: %v", err)
	}

	tests := []struct {
		key  string
		want string
	}{
		{"TEST_MULTI_A", "value_a"},
		{"TEST_MULTI_B", "value_b"},
		{"TEST_MULTI_C", "value_c"},
	}

	for _, tt := range tests {
		got := os.Getenv(tt.key)
		if got != tt.want {
			t.Errorf("%s = %q, want %q", tt.key, got, tt.want)
		}
	}

	// Cleanup
	os.Unsetenv("TEST_MULTI_A")
	os.Unsetenv("TEST_MULTI_B")
	os.Unsetenv("TEST_MULTI_C")
}

func TestLoadEnv_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	// Create a directory where the file should be (causes a read error)
	if err := os.Mkdir(envPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	err := loadEnv(envPath)
	if err == nil {
		t.Error("loadEnv() should return error for invalid .env (directory instead of file)")
	}
}
