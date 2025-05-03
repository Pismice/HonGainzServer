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

	// Handler to get all real_sessions of a user
	r.GET("/stats/real-sessions", func(c *gin.Context) {
		userID := c.GetInt("user_id")

		// Updated query to join with template_sessions and retrieve the name
		rows, err := db.Query(`
			SELECT rs.id, rs.template_session_id, rs.real_workout_id, rs.start_date, rs.finish_date, ts.name
			FROM real_sessions rs
			LEFT JOIN template_sessions ts ON rs.template_session_id = ts.id
			WHERE rs.user_id = ?`, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve real sessions"})
			return
		}

		defer rows.Close()

		// Updated to include template_session name
		var realSessions []map[string]interface{}
		for rows.Next() {
			var id, templateSessionID, realWorkoutID sql.NullInt64
			var startDate, finishDate, templateSessionName sql.NullString
			if err := rows.Scan(&id, &templateSessionID, &realWorkoutID, &startDate, &finishDate, &templateSessionName); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse real sessions"})
				return
			}
			realSessions = append(realSessions, map[string]interface{}{
				"id":                  id.Int64,
				"template_session_id": templateSessionID.Int64,
				"real_workout_id":     realWorkoutID.Int64,
				"start_date":          startDate.String,
				"finish_date":         finishDate.String,
				"name":                templateSessionName.String,
			})
		}

		// Return all real_sessions
		c.JSON(http.StatusOK, gin.H{
			"real_sessions": realSessions,
		})
	})
}
