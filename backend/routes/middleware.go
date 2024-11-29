package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func (routes *Routes) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins (use "*" for wide access or specify specific domains)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specific headers and methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If it's a preflight request, end here
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func (routes *Routes) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		fmt.Printf(authHeader)
		if authHeader == "" {
			data := Message{Data: "invalid or no jwt token"}
			sendJSONResponse(w, http.StatusBadRequest, data)
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer "){
			data := Message{Data: "missing bearer"}
			sendJSONResponse(w, http.StatusBadRequest, data)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Check if token is expired
		if claims.ExpiresAt < jwt.TimeFunc().Unix() {
			http.Error(w, "Unauthorized: Token expired", http.StatusUnauthorized)
			return
		}

		// Token is valid; proceed to the next handler
		next.ServeHTTP(w, r)
	})
}