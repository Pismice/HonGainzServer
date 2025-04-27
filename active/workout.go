package active

import (
	"database/sql"
	"fmt"
	"hongym/helper"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterActiveWorkoutRoutes(r *gin.RouterGroup, db *sql.DB) {
	// Fetch the active workout and its sessions
	r.GET("/active-workout", func(c *gin.Context) {
		userID := c.GetInt("user_id")

		// Fetch the active workout ID for the user
		var activeWorkoutID int
		err := db.QueryRow("SELECT active_workout_id FROM users WHERE id = ?", userID).Scan(&activeWorkoutID)
		if err != nil || activeWorkoutID == 0 {
			c.JSON(http.StatusOK, gin.H{"active_workout": nil})
			return
		}

		// Fetch the active workout details from the real_workouts table
		var startDate sql.NullString
		var templateWorkoutID int
		var templateWorkoutName string
		err = db.QueryRow(`
			SELECT rw.start_date, rw.template_workout_id, tw.name
			FROM real_workouts rw
			JOIN template_workouts tw ON rw.template_workout_id = tw.id
			WHERE rw.id = ?`, activeWorkoutID).Scan(&startDate, &templateWorkoutID, &templateWorkoutName)
		if err != nil {
			log.Printf("Error retrieving active workout details for workout ID %d: %v", activeWorkoutID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve active workout details"})
			return
		}

		// Fetch the sessions associated with the active workout from the real_sessions table
		sessionRows, err := db.Query(`
			SELECT rs.id, rs.start_date, rs.finish_date, ts.name
			FROM real_sessions rs
			JOIN template_sessions ts ON rs.template_session_id = ts.id
			WHERE rs.real_workout_id = ?`, activeWorkoutID)
		if err != nil {
			log.Printf("Error retrieving sessions for active workout ID %d: %v", activeWorkoutID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sessions for active workout"})
			return
		}
		defer sessionRows.Close()

		var sessions []map[string]interface{}
		for sessionRows.Next() {
			var sessionID int
			var sessionStartDate, sessionFinishDate sql.NullString
			var sessionName string
			if err := sessionRows.Scan(&sessionID, &sessionStartDate, &sessionFinishDate, &sessionName); err != nil {
				log.Printf("Error parsing session details for active workout ID %d: %v", activeWorkoutID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse session details"})
				return
			}

			// Convert nullable fields to valid strings or null
			sessions = append(sessions, map[string]interface{}{
				"id":          sessionID,
				"name":        sessionName,
				"start_date":  helper.NullStringToPointer(sessionStartDate),
				"finish_date": helper.NullStringToPointer(sessionFinishDate),
			})
		}

		// Return the active workout details with its sessions
		c.JSON(http.StatusOK, gin.H{
			"active_workout": map[string]interface{}{
				"id":                  activeWorkoutID,
				"start_date":          helper.NullStringToPointer(startDate),
				"template_workout_id": templateWorkoutID,
				"name":                templateWorkoutName,
				"user_id":             userID,
				"sessions":            sessions,
			},
		})
	})

	// Set the active workout for the connected user
	r.POST("/active-workout", func(c *gin.Context) {
		var payload struct {
			WorkoutID int `json:"workout_id" binding:"required"`
		}

		// Bind JSON input to the payload struct
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Invalid input: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		userID := c.GetInt("user_id")

		// Verify that the workout belongs to the user
		var ownerID int
		err := db.QueryRow("SELECT user_id FROM template_workouts WHERE id = ?", payload.WorkoutID).Scan(&ownerID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("Workout not found: %d", payload.WorkoutID)
				c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
			} else {
				log.Printf("Error verifying workout ownership for workout ID %d: %v", payload.WorkoutID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify workout ownership"})
			}
			return
		}

		if ownerID != userID {
			log.Printf("User ID %d does not have permission to activate workout ID %d", userID, payload.WorkoutID)
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to activate this workout"})
			return
		}

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		// Insert a new entry into the real_workouts table
		result, err := tx.Exec(`
			INSERT INTO real_workouts (start_date, template_workout_id, user_id)
			VALUES (?, ?, ?)`,
			time.Now().Format(time.RFC3339), payload.WorkoutID, userID)
		if err != nil {
			log.Printf("Error creating real workout: %v", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create real workout"})
			return
		}

		// Get the ID of the newly created real workout
		realWorkoutID, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error retrieving real workout ID: %v", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve real workout ID"})
			return
		}

		// Fetch all template sessions related to the template workout
		sessionRows, err := tx.Query(`
			SELECT template_session_id
			FROM template_workouts_template_sessions
			WHERE template_workout_id = ?`, payload.WorkoutID)
		if err != nil {
			log.Printf("Error retrieving template sessions for workout ID %d: %v", payload.WorkoutID, err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve template sessions"})
			return
		}
		defer sessionRows.Close()

		// Insert each template session into the real_sessions table
		for sessionRows.Next() {
			var templateSessionID int
			if err := sessionRows.Scan(&templateSessionID); err != nil {
				log.Printf("Error parsing template session ID for workout ID %d: %v", payload.WorkoutID, err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse template session ID"})
				return
			}

			// Insert the real session
			result, err := tx.Exec(`
				INSERT INTO real_sessions (template_session_id, real_workout_id, user_id)
				VALUES (?, ?, ?)`,
				templateSessionID, realWorkoutID, userID)
			if err != nil {
				log.Printf("Error creating real session for template session ID %d: %v", templateSessionID, err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create real session for template session ID " + fmt.Sprint(templateSessionID)})
				return
			}

			// Get the ID of the newly created real session
			realSessionID, err := result.LastInsertId()
			if err != nil {
				log.Printf("Error retrieving real session ID: %v", err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve real session ID"})
				return
			}

			// Fetch all template exercises corresponding to the template session
			exerciseRows, err := tx.Query(`
				SELECT template_exercise_id
				FROM template_sessions_template_exercises
				WHERE template_session_id = ?`, templateSessionID)
			if err != nil {
				log.Printf("Error retrieving template exercises for template session ID %d: %v", templateSessionID, err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve template exercises"})
				return
			}
			defer exerciseRows.Close()

			// Insert a real exercise for each template exercise
			for exerciseRows.Next() {
				var templateExerciseID int
				if err := exerciseRows.Scan(&templateExerciseID); err != nil {
					log.Printf("Error parsing template exercise ID for template session ID %d: %v", templateSessionID, err)
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse template exercise ID"})
					return
				}

				_, err := tx.Exec(`
					INSERT INTO real_exercises (template_exercise_id, real_session_id, user_id)
					VALUES (?, ?, ?)`,
					templateExerciseID, realSessionID, userID)
				if err != nil {
					log.Printf("Error creating real exercise for template exercise ID %d: %v", templateExerciseID, err)
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create real exercise for template exercise"})
					return
				}
			}
		}

		// Set the active workout for the user
		_, err = tx.Exec("UPDATE users SET active_workout_id = ? WHERE id = ?", realWorkoutID, userID)
		if err != nil {
			log.Printf("Error setting active workout for user ID %d: %v", userID, err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set active workout"})
			return
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Active workout set successfully with real exercises created"})
	})

	// Deactivate the active workout for the connected user
	r.DELETE("/active-workout", func(c *gin.Context) {
		userID := c.GetInt("user_id")

		// Retrieve the active workout ID for the user
		var activeWorkoutID int
		err := db.QueryRow("SELECT active_workout_id FROM users WHERE id = ?", userID).Scan(&activeWorkoutID)
		if err != nil || activeWorkoutID == 0 {
			log.Printf("No active workout to deactivate for user ID %d", userID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No active workout to deactivate"})
			return
		}

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		// Set the finish_date for the active workout
		_, err = tx.Exec("UPDATE real_workouts SET finish_date = ? WHERE id = ?", time.Now().Format(time.RFC3339), activeWorkoutID)
		if err != nil {
			log.Printf("Error setting finish date for active workout ID %d: %v", activeWorkoutID, err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set finish date for active workout"})
			return
		}

		// Set the active_workout_id field to NULL for the user
		_, err = tx.Exec("UPDATE users SET active_workout_id = NULL WHERE id = ?", userID)
		if err != nil {
			log.Printf("Error deactivating active workout for user ID %d: %v", userID, err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate active workout"})
			return
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Active workout deactivated successfully"})
	})

	// Add new real_sessions based on a real_workout
	r.POST("/new-sessions", func(c *gin.Context) {
		var payload struct {
			RealWorkoutID int `json:"real_workout_id" binding:"required"`
		}

		// Bind JSON input to the payload struct
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Invalid input: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		userID := c.GetInt("user_id")

		// Fetch the template_workout_id associated with the real_workout_id
		var templateWorkoutID int
		err := db.QueryRow("SELECT template_workout_id FROM real_workouts WHERE id = ? AND user_id = ?", payload.RealWorkoutID, userID).Scan(&templateWorkoutID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("Real workout not found or does not belong to user: %d", payload.RealWorkoutID)
				c.JSON(http.StatusNotFound, gin.H{"error": "Real workout not found or does not belong to user"})
			} else {
				log.Printf("Error retrieving template workout ID for real workout ID %d: %v", payload.RealWorkoutID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve template workout ID"})
			}
			return
		}

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		// Fetch all template sessions related to the template workout
		sessionRows, err := tx.Query(`
			SELECT template_session_id
			FROM template_workouts_template_sessions
			WHERE template_workout_id = ?`, templateWorkoutID)
		if err != nil {
			log.Printf("Error retrieving template sessions for template workout ID %d: %v", templateWorkoutID, err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve template sessions"})
			return
		}
		defer sessionRows.Close()

		// Insert each template session into the real_sessions table
		for sessionRows.Next() {
			var templateSessionID int
			if err := sessionRows.Scan(&templateSessionID); err != nil {
				log.Printf("Error parsing template session ID for template workout ID %d: %v", templateWorkoutID, err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse template session ID"})
				return
			}

			// Insert the real session
			result, err := tx.Exec(`
				INSERT INTO real_sessions (template_session_id, real_workout_id, user_id)
				VALUES (?, ?, ?)`,
				templateSessionID, payload.RealWorkoutID, userID)
			if err != nil {
				log.Printf("Error creating real session for template session ID %d: %v", templateSessionID, err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create real session for template session ID"})
				return
			}

			// Get the ID of the newly created real session
			realSessionID, err := result.LastInsertId()
			if err != nil {
				log.Printf("Error retrieving real session ID: %v", err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve real session ID"})
				return
			}

			// Fetch all template exercises corresponding to the template session
			exerciseRows, err := tx.Query(`
				SELECT template_exercise_id
				FROM template_sessions_template_exercises
				WHERE template_session_id = ?`, templateSessionID)
			if err != nil {
				log.Printf("Error retrieving template exercises for template session ID %d: %v", templateSessionID, err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve template exercises"})
				return
			}
			defer exerciseRows.Close()

			// Insert a real exercise for each template exercise
			for exerciseRows.Next() {
				var templateExerciseID int
				if err := exerciseRows.Scan(&templateExerciseID); err != nil {
					log.Printf("Error parsing template exercise ID for template session ID %d: %v", templateSessionID, err)
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse template exercise ID"})
					return
				}

				_, err := tx.Exec(`
					INSERT INTO real_exercises (template_exercise_id, real_session_id, user_id)
					VALUES (?, ?, ?)`,
					templateExerciseID, realSessionID, userID)
				if err != nil {
					log.Printf("Error creating real exercise for template exercise ID %d: %v", templateExerciseID, err)
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create real exercise for template exercise"})
					return
				}
			}
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Real sessions and exercises created successfully"})
	})
}
