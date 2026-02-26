package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey contextKey = "user"

// UserFromContext extracts the JWT claims from the request context.
// Returns nil if no authenticated user is present.
func UserFromContext(ctx context.Context) *Claims {
	claims, _ := ctx.Value(userContextKey).(*Claims)
	return claims
}

// Middleware returns an HTTP middleware that validates Bearer tokens.
// Requests without a valid token receive 401 Unauthorized.
func Middleware(jwtSvc *JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				slog.Warn("auth: missing Authorization header",
					slog.String("path", r.URL.Path),
				)
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				slog.Warn("auth: invalid Authorization format",
					slog.String("path", r.URL.Path),
				)
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtSvc.ValidateToken(tokenString)
			if err != nil {
				slog.Warn("auth: invalid token",
					slog.String("path", r.URL.Path),
					slog.String("error", err.Error()),
				)
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
