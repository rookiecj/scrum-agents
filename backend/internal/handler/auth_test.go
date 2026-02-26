package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rookiecj/scrum-agents/backend/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func testStore(t *testing.T) *auth.Store {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	store, err := auth.NewStore(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestHandleSignup(t *testing.T) {
	store := testStore(t)
	handler := HandleSignup(store)

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "valid signup",
			body:       `{"email":"alice@example.com","password":"password123"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid email",
			body:       `{"email":"not-an-email","password":"password123"}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "short password",
			body:       `{"email":"bob@example.com","password":"short"}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "invalid json",
			body:       `not json`,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "empty body",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/signup", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", rec.Code, tt.wantStatus, rec.Body.String())
			}

			var resp SignupResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if tt.wantErr && resp.Error == "" {
				t.Error("expected error in response")
			}
			if !tt.wantErr && resp.Error != "" {
				t.Errorf("unexpected error: %s", resp.Error)
			}
		})
	}
}

func TestHandleSignup_DuplicateEmail(t *testing.T) {
	store := testStore(t)
	handler := HandleSignup(store)

	// First signup
	body := `{"email":"alice@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("first signup: status = %d, want %d", rec.Code, http.StatusCreated)
	}

	// Duplicate signup
	req = httptest.NewRequest("POST", "/api/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("duplicate signup: status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestHandleLogin(t *testing.T) {
	store := testStore(t)
	jwtSvc := auth.NewJWTService("test-secret", time.Hour)

	// Pre-create a user with bcrypt-hashed password
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if _, err := store.CreateUser("alice@example.com", string(hash)); err != nil {
		t.Fatal(err)
	}

	handler := HandleLogin(store, jwtSvc)

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantToken  bool
		wantErr    bool
	}{
		{
			name:       "valid login",
			body:       `{"email":"alice@example.com","password":"password123"}`,
			wantStatus: http.StatusOK,
			wantToken:  true,
		},
		{
			name:       "wrong password",
			body:       `{"email":"alice@example.com","password":"wrongpass123"}`,
			wantStatus: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "non-existent user",
			body:       `{"email":"nobody@example.com","password":"password123"}`,
			wantStatus: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "invalid json",
			body:       `not json`,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", rec.Code, tt.wantStatus, rec.Body.String())
			}

			var resp LoginResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if tt.wantToken && resp.Token == "" {
				t.Error("expected token in response")
			}
			if tt.wantErr && resp.Error == "" {
				t.Error("expected error in response")
			}
		})
	}
}

func TestHandleLogin_TokenIsValid(t *testing.T) {
	store := testStore(t)
	jwtSvc := auth.NewJWTService("test-secret", time.Hour)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if _, err := store.CreateUser("alice@example.com", string(hash)); err != nil {
		t.Fatal(err)
	}

	handler := HandleLogin(store, jwtSvc)

	body := `{"email":"alice@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var resp LoginResponse
	json.NewDecoder(rec.Body).Decode(&resp)

	// Validate the returned token
	claims, err := jwtSvc.ValidateToken(resp.Token)
	if err != nil {
		t.Fatalf("returned token is invalid: %v", err)
	}
	if claims.Email != "alice@example.com" {
		t.Errorf("token email = %q, want %q", claims.Email, "alice@example.com")
	}
}
