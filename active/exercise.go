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

	// Get the number of real_sets for a given real_exercise_id
	r.GET("/real-sets/count/:real_exercise_id", func(c *gin.Context) {
		realExerciseID := c.Param("real_exercise_id")
		userID := c.GetInt("user_id")

		// Query to count the number of real_sets
		var count int
		err := db.QueryRow(`
			SELECT COUNT(*)
			FROM real_sets
			WHERE real_exercise_id = ? AND user_id = ?`,
			realExerciseID, userID).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve count of real_sets"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"real_exercise_id": realExerciseID,
			"real_sets_count":  count,
		})
	})
}
