package auth

import (
	"testing"
	"time"
)

func TestGenerateAndValidateToken(t *testing.T) {
	svc := NewJWTService("test-secret", time.Hour)

	token, err := svc.GenerateToken(42, "alice@example.com")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}
	if claims.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", claims.Email, "alice@example.com")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	svc := NewJWTService("test-secret", -time.Hour) // already expired

	token, err := svc.GenerateToken(1, "bob@example.com")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, err = svc.ValidateToken(token)
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	svc1 := NewJWTService("secret-1", time.Hour)
	svc2 := NewJWTService("secret-2", time.Hour)

	token, _ := svc1.GenerateToken(1, "test@example.com")

	_, err := svc2.ValidateToken(token)
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	svc := NewJWTService("test-secret", time.Hour)

	_, err := svc.ValidateToken("not-a-valid-token")
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}
