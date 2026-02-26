package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"

	"github.com/rookiecj/scrum-agents/backend/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	ID    int64  `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Error string `json:"error,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// HandleSignup returns a handler for POST /api/signup.
func HandleSignup(store *auth.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Warn("signup: invalid request body",
				slog.String("handler", "signup"),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusBadRequest, SignupResponse{Error: "invalid request body"})
			return
		}

		req.Email = strings.TrimSpace(strings.ToLower(req.Email))

		// Validate email format
		if _, err := mail.ParseAddress(req.Email); err != nil {
			writeJSON(w, http.StatusBadRequest, SignupResponse{Error: "invalid email format"})
			return
		}

		// Validate password length
		if len(req.Password) < 8 {
			writeJSON(w, http.StatusBadRequest, SignupResponse{Error: "password must be at least 8 characters"})
			return
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("signup: password hashing failed",
				slog.String("handler", "signup"),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusInternalServerError, SignupResponse{Error: "internal server error"})
			return
		}

		user, err := store.CreateUser(req.Email, string(hash))
		if err != nil {
			if errors.Is(err, auth.ErrUserExists) {
				writeJSON(w, http.StatusConflict, SignupResponse{Error: "email already registered"})
				return
			}
			slog.Error("signup: create user failed",
				slog.String("handler", "signup"),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusInternalServerError, SignupResponse{Error: "internal server error"})
			return
		}

		slog.Info("signup: user created",
			slog.String("handler", "signup"),
			slog.Int64("user_id", user.ID),
			slog.String("email", user.Email),
		)
		writeJSON(w, http.StatusCreated, SignupResponse{
			ID:    user.ID,
			Email: user.Email,
		})
	}
}

// HandleLogin returns a handler for POST /api/login.
func HandleLogin(store *auth.Store, jwtSvc *auth.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Warn("login: invalid request body",
				slog.String("handler", "login"),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusBadRequest, LoginResponse{Error: "invalid request body"})
			return
		}

		req.Email = strings.TrimSpace(strings.ToLower(req.Email))

		user, err := store.GetUserByEmail(req.Email)
		if err != nil {
			if errors.Is(err, auth.ErrUserNotFound) {
				writeJSON(w, http.StatusUnauthorized, LoginResponse{Error: "invalid email or password"})
				return
			}
			slog.Error("login: query user failed",
				slog.String("handler", "login"),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusInternalServerError, LoginResponse{Error: "internal server error"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			writeJSON(w, http.StatusUnauthorized, LoginResponse{Error: "invalid email or password"})
			return
		}

		token, err := jwtSvc.GenerateToken(user.ID, user.Email)
		if err != nil {
			slog.Error("login: token generation failed",
				slog.String("handler", "login"),
				slog.String("error", err.Error()),
			)
			writeJSON(w, http.StatusInternalServerError, LoginResponse{Error: "internal server error"})
			return
		}

		slog.Info("login: success",
			slog.String("handler", "login"),
			slog.Int64("user_id", user.ID),
		)
		writeJSON(w, http.StatusOK, LoginResponse{Token: token})
	}
}
