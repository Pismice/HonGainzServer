package stats

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterStatsRoutes(r *gin.RouterGroup, db *sql.DB) {
	// Get the biggest weight ever recorded for a given template_exercise
	r.GET("/stats/max-weight/:template_exercise_id", func(c *gin.Context) {
		templateExerciseID := c.Param("template_exercise_id")
		userID := c.GetInt("user_id")

		// Query to find the maximum weight for the given template_exercise
		var maxWeight sql.NullInt64
		err := db.QueryRow(`
			SELECT MAX(rs.weight)
			FROM real_sets rs
			JOIN real_exercises re ON rs.real_exercise_id = re.id
			WHERE re.template_exercise_id = ? AND rs.user_id = ?`, templateExerciseID, userID).Scan(&maxWeight)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve max weight"})
			return
		}

		// Return the maximum weight or null if no records exist
		c.JSON(http.StatusOK, gin.H{
			"template_exercise_id": templateExerciseID,
			"max_weight":           maxWeight.Int64,
		})
	})

	r.GET("/stats/all-weights/:template_exercise_id", func(c *gin.Context) {
		templateExerciseID := c.Param("template_exercise_id")
		userID := c.GetInt("user_id")

		// Query to find all weights for the given template_exercise
		rows, err := db.Query(`
			SELECT rs.weight
			FROM real_sets rs
			JOIN real_exercises re ON rs.real_exercise_id = re.id
			WHERE re.template_exercise_id = ? AND rs.user_id = ?`, templateExerciseID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve weights"})
			return
		}
		defer rows.Close()

		// Collect all weights
		var weights []int64
		for rows.Next() {
			var weight int64
			if err := rows.Scan(&weight); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse weights"})
				return
			}
			weights = append(weights, weight)
		}

		// Return all weights
		c.JSON(http.StatusOK, gin.H{
			"template_exercise_id": templateExerciseID,
			"weights":              weights,
		})
	})
}
