package auth

import (
	"os"
	"testing"
)

func tempDB(t *testing.T) *Store {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	store, err := NewStore(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestCreateUser(t *testing.T) {
	store := tempDB(t)

	user, err := store.CreateUser("alice@example.com", "hashed-password")
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if user.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if user.Email != "alice@example.com" {
		t.Errorf("email = %q, want %q", user.Email, "alice@example.com")
	}
	if user.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	store := tempDB(t)

	if _, err := store.CreateUser("bob@example.com", "hash1"); err != nil {
		t.Fatalf("first CreateUser() error = %v", err)
	}

	_, err := store.CreateUser("bob@example.com", "hash2")
	if err != ErrUserExists {
		t.Errorf("expected ErrUserExists, got %v", err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	store := tempDB(t)

	_, err := store.CreateUser("carol@example.com", "myhash")
	if err != nil {
		t.Fatal(err)
	}

	user, err := store.GetUserByEmail("carol@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail() error = %v", err)
	}

	if user.Email != "carol@example.com" {
		t.Errorf("email = %q, want %q", user.Email, "carol@example.com")
	}
	if user.PasswordHash != "myhash" {
		t.Errorf("password_hash = %q, want %q", user.PasswordHash, "myhash")
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	store := tempDB(t)

	_, err := store.GetUserByEmail("nobody@example.com")
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
