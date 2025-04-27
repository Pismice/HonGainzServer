package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OwnershipMiddleware(db *sql.DB, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		resourceID := c.Param("id")
		username := c.GetString("username")

		var userID int
		err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			c.Abort()
			return
		}

		var ownerID int
		query := "SELECT user_id FROM " + resourceType + " WHERE id = ?"
		err = db.QueryRow(query, resourceID).Scan(&ownerID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
			c.Abort()
			return
		}

		if ownerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
