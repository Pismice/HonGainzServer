package active

import (
	"database/sql"
	"hongym/helper"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterActiveSessionRoutes(r *gin.RouterGroup, db *sql.DB) {
	// Fetch the active session details
	r.GET("/active-session/:id", func(c *gin.Context) {
		sessionID := c.Param("id")
		userID := c.GetInt("user_id")

		// Fetch the session details from the real_sessions table
		var startDate, finishDate sql.NullString
		var templateSessionID int
		var sessionName string
		err := db.QueryRow(`
			SELECT rs.start_date, rs.finish_date, rs.template_session_id, ts.name
			FROM real_sessions rs
			JOIN template_sessions ts ON rs.template_session_id = ts.id
			WHERE rs.id = ? AND rs.user_id = ?`, sessionID, userID).Scan(&startDate, &finishDate, &templateSessionID, &sessionName)
		if err != nil {
			log.Printf("Error retrieving session details for session ID %s: %v", sessionID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session details"})
			return
		}

		// Return the session details
		c.JSON(http.StatusOK, gin.H{
			"session": map[string]interface{}{
				"id":                  sessionID,
				"start_date":          helper.NullStringToPointer(startDate),
				"finish_date":         helper.NullStringToPointer(finishDate),
				"template_session_id": templateSessionID,
				"name":                sessionName,
			},
		})
	})

	// Start a session
	r.POST("/active-session/:id/start", func(c *gin.Context) {
		sessionID := c.Param("id")
		userID := c.GetInt("user_id")

		// Start the session by setting the start_date
		_, err := db.Exec("UPDATE real_sessions SET start_date = ? WHERE id = ? AND user_id = ?", time.Now().Format(time.RFC3339), sessionID, userID)
		if err != nil {
			log.Printf("Error starting session ID %s: %v", sessionID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Session started successfully"})
	})

	// Finish a session
	r.POST("/active-session/:id/finish", func(c *gin.Context) {
		sessionID := c.Param("id")
		userID := c.GetInt("user_id")

		// Finish the session by setting the finish_date
		_, err := db.Exec("UPDATE real_sessions SET finish_date = ? WHERE id = ? AND user_id = ?", time.Now().Format(time.RFC3339), sessionID, userID)
		if err != nil {
			log.Printf("Error finishing session ID %s: %v", sessionID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Session finished successfully"})
	})

	// Get all real exercises for a real session
	r.GET("/real-session/:id/exercises", func(c *gin.Context) {
		sessionID := c.Param("id")
		userID := c.GetInt("user_id")

		// Query to fetch all real exercises with their names, start_date, and finish_date for the given real session
		rows, err := db.Query(`
			SELECT re.id, re.template_exercise_id, te.name, re.start_date, re.finish_date
			FROM real_exercises re
			JOIN template_exercises te ON re.template_exercise_id = te.id
			WHERE re.real_session_id = ? AND re.user_id = ?`, sessionID, userID)
		if err != nil {
			log.Printf("Error retrieving real exercises for session ID %s: %v", sessionID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve real exercises"})
			return
		}
		defer rows.Close()

		// Parse the results
		var exercises []map[string]interface{}
		for rows.Next() {
			var exerciseID, templateExerciseID int
			var exerciseName string
			var startDate, finishDate sql.NullString
			if err := rows.Scan(&exerciseID, &templateExerciseID, &exerciseName, &startDate, &finishDate); err != nil {
				log.Printf("Error parsing real exercise for session ID %s: %v", sessionID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse real exercises"})
				return
			}

			exercises = append(exercises, map[string]interface{}{
				"id":                   exerciseID,
				"template_exercise_id": templateExerciseID,
				"name":                 exerciseName,
				"start_date":           helper.NullStringToPointer(startDate),
				"finish_date":          helper.NullStringToPointer(finishDate),
			})
		}

		// Ensure the response is valid JSON
		if exercises == nil {
			exercises = []map[string]interface{}{}
		}

		// Return the list of real exercises
		c.JSON(http.StatusOK, gin.H{"real_exercises": exercises})
	})

}
