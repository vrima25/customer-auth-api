package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"ticket-triage-api/util"
)

type contextKey string

const UserContextKey contextKey = "user"

func writeJSONError(w http.ResponseWriter, status int, message string){
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer"){
				writeJSONError(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return 
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := util.ParseToken(tokenString, jwtSecret)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "invalid or expired token")
				return 
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}