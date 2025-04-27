package middleware

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		sessionID := strings.TrimPrefix(authHeader, "Bearer ")

		var username string
		var userID int
		// Retrieve both username and user ID from the database
		err := db.QueryRow("SELECT id, username FROM users WHERE session_id = ?", sessionID).Scan(&userID, &username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		// Set the username and user ID in the Gin context
		c.Set("username", username)
		c.Set("user_id", userID)
		c.Next()
	}
}
