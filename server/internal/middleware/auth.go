package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("Content-Type", "application/json")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - No token provided",
			})
			return
		}

		// Check if the header starts with "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Unauthorized - Invalid token format",
			})
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Unauthorized - Invalid token",
			})
			return
		}

		// Get claims from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Unauthorized - Invalid token claims",
			})
			return
		}

		// Convert user ID to float64
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Unauthorized - Invalid user ID",
			})
			return
		}
		c.Set("user_id", userID)
		c.Set("role", claims["role"])

		c.Next()
	}
}
