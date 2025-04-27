package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func generateSessionID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal("Failed to generate session ID:", err)
	}
	return hex.EncodeToString(bytes)
}

func RegisterAuthRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/login", func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", credentials.Username).Scan(&storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		if storedPassword != credentials.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		sessionID := generateSessionID()
		_, err = db.Exec("UPDATE users SET session_id = ? WHERE username = ?", sessionID, credentials.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session_id": sessionID})
	})

	r.POST("/register", func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var existingUser string
		err := db.QueryRow("SELECT username FROM users WHERE username = ?", credentials.Username).Scan(&existingUser)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		} else if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", credentials.Username, credentials.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		sessionID := generateSessionID()
		_, err = db.Exec("UPDATE users SET session_id = ? WHERE username = ?", sessionID, credentials.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "session_id": sessionID})
	})
}
