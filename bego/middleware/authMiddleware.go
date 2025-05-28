package middleware

import (
	"bego/config"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type key string

const (
	ContextUserID key = "user_id"
	ContextRole   key = "role"
)

// ProtectRoutes is middleware to require valid JWT token
func ProtectRoutes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		//Extract the token from header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		//Parse and verify the token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JwtSecretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized - invalid token", http.StatusUnauthorized)
			return
		}

		//Extract values from claims and add to request context
		userID, ok1 := claims["user_id"].(string)
		role, ok2 := claims["role"].(string)
		if !ok1 || !ok2 {
			http.Error(w, "Unauthorized - invalid claims", http.StatusUnauthorized)
			return
		}

		//Add user_id and role to context
		ctx := context.WithValue(r.Context(), ContextUserID, userID)
		ctx = context.WithValue(ctx, ContextRole, role)

		//Pass modified request to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
