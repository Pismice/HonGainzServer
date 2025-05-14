package active

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterExerciseRoutes(r *gin.RouterGroup, db *sql.DB) {
	// Register a new set for a real exercise
	r.POST("/register-set", func(c *gin.Context) {
		var payload struct {
			RealExerciseID int     `json:"real_exercise_id" binding:"required"`
			Reps           int     `json:"reps" binding:"required"`
			Weight         float64 `json:"weight" binding:"required"`
		}

		// Bind JSON input to the payload struct
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		userID := c.GetInt("user_id")

		// Insert the new set into the real_sets table
		_, err := db.Exec(`
			INSERT INTO real_sets (reps, weight, real_exercise_id, user_id, start_date, finish_date)
			VALUES (?, ?, ?, ?, ?, ?)`,
			payload.Reps, payload.Weight, payload.RealExerciseID, userID, time.Now().Format(time.RFC3339), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register set"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Set registered successfully"})
	})

	// Set a finish date for a real exercise
	r.POST("/finish-exercise/:id", func(c *gin.Context) {
		exerciseID := c.Param("id")
		userID := c.GetInt("user_id")

		// Update the finish_date for the real exercise
		_, err := db.Exec(`
			UPDATE real_exercises
			SET finish_date = ?
			WHERE id = ? AND user_id = ?`,
			time.Now().Format(time.RFC3339), exerciseID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set finish date for exercise"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Exercise finished successfully"})
	})
}
