package middleware

import (
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"fmt"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("Auth header:", authHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			return
		}
		splitToken := strings.Split(authHeader, "Bearer ")
		fmt.Println("Split token:", splitToken)
		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is invalid"})
			return
		}
		tokenString := splitToken[1]
		fmt.Println("Token string length:", len(tokenString))
		fmt.Println("Token string:", tokenString)
		
		// Trim any whitespace from the token
		tokenString = strings.TrimSpace(tokenString)
		
		// Check if token is properly formatted (should have 3 parts separated by dots)
		parts := strings.Split(tokenString, ".")
		if len(parts) != 3 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token format is invalid - should have 3 parts"})
			return
		}
		
		// Basic validation of each part
		for i, part := range parts {
			if part == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Token part %d is empty", i+1)})
				return
			}
		}
		
		userId, err := utils.ParseToken(tokenString)
		fmt.Println("userId", userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}
