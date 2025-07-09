package api

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/cliffdoyle/go-auth-app/internal/auth"
	"github.com/cliffdoyle/go-auth-app/internal/model"
	"net/http"
	"strings"
)

type contextKey string
const UserContextKey = contextKey("user")

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := auth.ValidateJWT(tokenString, jwtSecret)
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContextKey).(jwt.MapClaims)
		if !ok {
			http.Error(w, "No user data in context", http.StatusInternalServerError)
			return
		}
		role, ok := claims["role"].(string)
		if !ok || model.Role(role) != model.AdminRole {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}