package routes

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(email string) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	return token.SignedString(jwtSecret)
}