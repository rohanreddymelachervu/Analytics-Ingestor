package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware verifies JWT and stores claims in context
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing or invalid"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}
		// store claims in context for downstream handlers
		c.Set("userID", claims.UserID)
		c.Set("scopes", claims.Scopes)
		c.Next()
	}
}

// RequireScope checks that the token has the required scope
func RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopesIfc, exists := c.Get("scopes")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no scopes present"})
			return
		}
		scopes := scopesIfc.([]string)
		for _, s := range scopes {
			if s == scope {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient scope"})
	}
}
