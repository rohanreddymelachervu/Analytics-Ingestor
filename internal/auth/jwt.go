package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a token with embedded scopes (e.g. "READ", "WRITE")
func GenerateJWT(userID uint, secret string, scopes []string) (string, error) {
	claims := Claims{
		UserID: userID,
		Scopes: scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Claims extends standard JWT claims with user ID and scopes
type Claims struct {
	UserID uint     `json:"user_id"`
	Scopes []string `json:"scopes"`
	jwt.RegisteredClaims
}
