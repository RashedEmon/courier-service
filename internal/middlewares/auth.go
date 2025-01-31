package middleware

import (
	"courier-service/internal/helpers"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// this middleware that checks JWT authentication for /api/v1/orders routes.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse Token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		tokenString := tokenParts[1]

		// Validate token
		claims, err := helpers.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Pass user id and email to Gin context
		c.Set("user_id", claims.UserId)
		c.Set("email", claims.Email)

		fmt.Println("Request come from email: ", claims.Email, "User Id: ", claims.UserId)

		c.Next()
	}
}
