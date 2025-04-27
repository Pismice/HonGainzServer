package template

import (
	"database/sql"
	"net/http"

	"hongym/middleware"

	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterTemplateRoutes(r *gin.RouterGroup, db *sql.DB) {

	// Create a template workout with associated sessions
	r.POST("/template-workouts", func(c *gin.Context) {
		var payload struct {
			Name       string `json:"name" binding:"required"`
			SessionIDs []int  `json:"session_ids" binding:"required"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		userID := c.GetInt("user_id")

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

			// Insert the new workout into the template_workouts table
		result, err := tx.Exec("INSERT INTO template_workouts (name, user_id) VALUES (?, ?)", payload.Name, userID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template workout"})
			return
		}

		// Get the ID of the newly created workout
		workoutID, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workout ID"})
			return
		}

		// Insert the relationships into the template_workouts_template_sessions table
		for _, sessionID := range payload.SessionIDs {
			_, err := tx.Exec("INSERT INTO template_workouts_template_sessions (template_workout_id, template_session_id) VALUES (?, ?)", workoutID, sessionID)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate workout with sessions for session ID " + fmt.Sprint(sessionID)})
				return
			}
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Template workout created successfully"})
	})

	// Retrieve all template workouts with session IDs
	r.GET("/template-workouts", func(c *gin.Context) {
		userID := c.GetInt("user_id")

		rows, err := db.Query("SELECT id, name FROM template_workouts WHERE user_id = ?", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve template workouts"})
			return
		}
		defer rows.Close()

		var workouts []map[string]interface{}
		for rows.Next() {
			var workoutID int
			var name string
			if err := rows.Scan(&workoutID, &name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse template workouts"})
				return
			}

			// Retrieve session IDs for the workout
			sessionRows, err := db.Query("SELECT template_session_id FROM template_workouts_template_sessions WHERE template_workout_id = ?", workoutID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sessions for workout"})
				return
			}
			defer sessionRows.Close()

			var sessionIDs []int
			for sessionRows.Next() {
				var sessionID int
				if err := sessionRows.Scan(&sessionID); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse session IDs"})
					return
				}
				sessionIDs = append(sessionIDs, sessionID)
			}

			// Add workout with session IDs to the response
			workouts = append(workouts, map[string]interface{}{
				"id":          workoutID,
				"name":        name,
				"session_ids": sessionIDs,
			})
		}

		c.JSON(http.StatusOK, gin.H{"template_workouts": workouts})
	})

	// Update a template workout with associated sessions
	r.PUT("/template-workouts/:id", middleware.OwnershipMiddleware(db, "template_workouts"), func(c *gin.Context) {
		var payload struct {
			Name       string `json:"name" binding:"required"`
			SessionIDs []int  `json:"session_ids" binding:"required"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		workoutID := c.Param("id")

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		// Update the workout name in the template_workouts table
		_, err = tx.Exec("UPDATE template_workouts SET name = ? WHERE id = ?", payload.Name, workoutID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template workout name"})
			return
		}

		// Delete existing relationships in the template_workouts_template_sessions table
		_, err = tx.Exec("DELETE FROM template_workouts_template_sessions WHERE template_workout_id = ?", workoutID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear existing session relationships"})
			return
		}

		// Insert the new relationships into the template_workouts_template_sessions table
		for _, sessionID := range payload.SessionIDs {
			_, err := tx.Exec("INSERT INTO template_workouts_template_sessions (template_workout_id, template_session_id) VALUES (?, ?)", workoutID, sessionID)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate workout with sessions for session ID " + fmt.Sprint(sessionID)})
				return
			}
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Template workout updated successfully"})
	})

	// Delete a template workout
	r.DELETE("/template-workouts/:id", middleware.OwnershipMiddleware(db, "template_workouts"), func(c *gin.Context) {
		workoutID := c.Param("id")
		_, err := db.Exec("DELETE FROM template_workouts WHERE id = ?", workoutID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template workout"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Template workout deleted successfully"})
	})

	// Create a template session
	r.POST("/template-sessions", func(c *gin.Context) {
		var payload struct {
			Name        string `json:"name" binding:"required"`
			ExerciseIDs []int  `json:"exercise_ids" binding:"required"`
		}

		// Bind JSON input to the payload struct
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Retrieve the user ID from the context
		userID := c.GetInt("user_id")

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		// Insert the new session into the template_sessions table
		result, err := tx.Exec("INSERT INTO template_sessions (name, user_id) VALUES (?, ?)", payload.Name, userID)
		if err != nil {
			tx.Rollback() // Rollback the transaction on error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template session"})
			return
		}

		// Get the ID of the newly created session
		sessionID, err := result.LastInsertId()
		if err != nil {
			tx.Rollback() // Rollback the transaction on error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session ID"})
			return
		}

		// Insert the relationships into the template_sessions-template_exercises table
		for _, exerciseID := range payload.ExerciseIDs {
			_, err := tx.Exec("INSERT INTO template_sessions_template_exercises (template_session_id, template_exercise_id) VALUES (?, ?)", sessionID, exerciseID)
			if err != nil {
				tx.Rollback() // Rollback the transaction on error
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate session with exercises for exercise ID " + fmt.Sprint(exerciseID) + " and session id " + fmt.Sprint(sessionID)})
				return
			}
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			tx.Rollback() // Rollback the transaction on error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "Template session created successfully"})
	})

	// Retrieve all template sessions with exercise IDs
	r.GET("/template-sessions", func(c *gin.Context) {
		userID := c.GetInt("user_id")

		rows, err := db.Query("SELECT id, name FROM template_sessions WHERE user_id = ?", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve template sessions"})
			return
		}
		defer rows.Close()

		var sessions []map[string]interface{}
		for rows.Next() {
			var sessionID int
			var name string
			if err := rows.Scan(&sessionID, &name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse template sessions"})
				return
			}

			// Retrieve exercise IDs for the session
			exerciseRows, err := db.Query("SELECT template_exercise_id FROM template_sessions_template_exercises WHERE template_session_id = ?", sessionID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exercises for session"})
				return
			}
			defer exerciseRows.Close()

			var exerciseIDs []int
			for exerciseRows.Next() {
				var exerciseID int
				if err := exerciseRows.Scan(&exerciseID); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse exercise IDs"})
					return
				}
				exerciseIDs = append(exerciseIDs, exerciseID)
			}

			// Add session with exercise IDs to the response
			sessions = append(sessions, map[string]interface{}{
				"id":           sessionID,
				"name":         name,
				"exercise_ids": exerciseIDs,
			})
		}

		c.JSON(http.StatusOK, gin.H{"template_sessions": sessions})
	})

	// Update a template session
	r.PUT("/template-sessions/:id", middleware.OwnershipMiddleware(db, "template_sessions"), func(c *gin.Context) {
		var payload struct {
			Name        string `json:"name" binding:"required"`
			ExerciseIDs []int  `json:"exercise_ids" binding:"required"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		sessionID := c.Param("id")

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		// Update the session name in the template_sessions table
		_, err = tx.Exec("UPDATE template_sessions SET name = ? WHERE id = ?", payload.Name, sessionID)
		if err != nil {
			tx.Rollback() // Rollback the transaction on error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template session name"})
			return
		}

		// Delete existing relationships in the template_sessions_template_exercises table
		_, err = tx.Exec("DELETE FROM template_sessions_template_exercises WHERE template_session_id = ?", sessionID)
		if err != nil {
			tx.Rollback() // Rollback the transaction on error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear existing exercise relationships"})
			return
		}

		// Insert the new relationships into the template_sessions_template_exercises table
		for _, exerciseID := range payload.ExerciseIDs {
			_, err := tx.Exec("INSERT INTO template_sessions_template_exercises (template_session_id, template_exercise_id) VALUES (?, ?)", sessionID, exerciseID)
			if err != nil {
				tx.Rollback() // Rollback the transaction on error
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate session with exercises for exercise ID " + fmt.Sprint(exerciseID)})
				return
			}
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			tx.Rollback() // Rollback the transaction on error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "Template session updated successfully"})
	})

	// Delete a template session
	r.DELETE("/template-sessions/:id", middleware.OwnershipMiddleware(db, "template_sessions"), func(c *gin.Context) {
		sessionID := c.Param("id")
		_, err := db.Exec("DELETE FROM template_sessions WHERE id = ?", sessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Template session deleted successfully"})
	})

	// Create a template exercise
	r.POST("/template-exercises", func(c *gin.Context) {
		var payload struct {
			Name string `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		userID := c.GetInt("user_id")

		_, err := db.Exec("INSERT INTO template_exercises (name, user_id) VALUES (?, ?)", payload.Name, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template exercise"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Template exercise created successfully"})
	})

	// Retrieve all template exercises
	r.GET("/template-exercises", func(c *gin.Context) {
		userID := c.GetInt("user_id")

		rows, err := db.Query("SELECT id, name FROM template_exercises WHERE user_id = ?", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve template exercises"})
			return
		}
		defer rows.Close()

		var exercises []map[string]interface{}
		for rows.Next() {
			var id int
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse template exercises"})
				return
			}
			exercises = append(exercises, map[string]interface{}{"id": id, "name": name})
		}

		c.JSON(http.StatusOK, gin.H{"template_exercises": exercises})
	})

	// Update a template exercise
	r.PUT("/template-exercises/:id", middleware.OwnershipMiddleware(db, "template_exercises"), func(c *gin.Context) {
		var payload struct {
			Name string `json:"name" binding:"required"`
		}

		println("TRYING TO PUT")

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		exerciseID := c.Param("id")
		_, err := db.Exec("UPDATE template_exercises SET name = ? WHERE id = ?", payload.Name, exerciseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template exercise"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Template exercise updated successfully"})
	})

	// Delete a template exercise
	r.DELETE("/template-exercises/:id", middleware.OwnershipMiddleware(db, "template_exercises"), func(c *gin.Context) {
		println("TRYING TO DELETE")
		exerciseID := c.Param("id")
		_, err := db.Exec("DELETE FROM template_exercises WHERE id = ?", exerciseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template exercise"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Template exercise deleted successfully"})
	})
}
