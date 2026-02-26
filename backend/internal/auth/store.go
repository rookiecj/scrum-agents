package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/rookiecj/scrum-agents/backend/internal/model"

	_ "modernc.org/sqlite"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

// Store manages user persistence with SQLite.
type Store struct {
	db *sql.DB
}

// NewStore opens (or creates) a SQLite database at the given path and
// initialises the users table.
func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	const createTable = `
		CREATE TABLE IF NOT EXISTS users (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			email         TEXT    NOT NULL UNIQUE,
			password_hash TEXT    NOT NULL,
			created_at    TEXT    NOT NULL DEFAULT (datetime('now'))
		);`

	if _, err := db.Exec(createTable); err != nil {
		return nil, fmt.Errorf("create users table: %w", err)
	}

	return &Store{db: db}, nil
}

// CreateUser inserts a new user. Returns ErrUserExists if the email is taken.
func (s *Store) CreateUser(email, passwordHash string) (*model.User, error) {
	const q = `INSERT INTO users (email, password_hash) VALUES (?, ?)`
	result, err := s.db.Exec(q, email, passwordHash)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}

	id, _ := result.LastInsertId()
	return &model.User{
		ID:        id,
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// GetUserByEmail returns a user by email. Returns ErrUserNotFound if absent.
func (s *Store) GetUserByEmail(email string) (*model.User, error) {
	const q = `SELECT id, email, password_hash, created_at FROM users WHERE email = ?`
	row := s.db.QueryRow(q, email)

	var u model.User
	var createdAt string
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query user: %w", err)
	}

	u.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return &u, nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// isUniqueViolation checks if the error is a SQLite UNIQUE constraint violation.
func isUniqueViolation(err error) bool {
	return err != nil && (errors.Is(err, sql.ErrNoRows) == false) &&
		(containsString(err.Error(), "UNIQUE constraint failed") ||
			containsString(err.Error(), "constraint failed"))
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
