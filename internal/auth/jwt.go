package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT creates a token with embedded scopes (e.g. "READ", "WRITE")
func GenerateJWT(userID uint, secret string, scopes []string) (string, error) {
	claims := Claims{
		UserID: userID,
		Scopes: scopes,
		StandardClaims: jwt.StandardClaims{
			Subject:   fmt.Sprint(userID),
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Claims extends standard JWT claims with user ID and scopes
type Claims struct {
	UserID uint     `json:"user_id"`
	Scopes []string `json:"scopes"`
	jwt.StandardClaims
}
