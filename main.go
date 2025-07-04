package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"hongym/active"
	"hongym/auth"
	"hongym/middleware"
	"hongym/stats" // Added import for stats
	"hongym/template"
)

// Set useHTTPS based on the environment variable "HON_GYM_PROD"
var useHTTPS = os.Getenv("HON_GYM_PROD") != ""

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()
	r.SetTrustedProxies(nil) // Désactive la confiance envers tout proxy


	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	// Add a middleware to log errors
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("[ERROR] %s", e.Error())
			}
		}
	})

	// Register authentication routes
	auth.RegisterAuthRoutes(r, db)

	// Authenticated routes
	authenticated := r.Group("/auth")
	authenticated.Use(middleware.AuthMiddleware(db))
	{
		template.RegisterTemplateRoutes(authenticated, db)    // Register template routes
		active.RegisterActiveSessionRoutes(authenticated, db) // Register active routes
		active.RegisterActiveWorkoutRoutes(authenticated, db) // Register active routes
		active.RegisterExerciseRoutes(authenticated, db)      // Register active routes
		stats.RegisterStatsRoutes(authenticated, db)          // Register active routes
		authenticated.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Everything is OK!"})
		})
	}

	r.GET("/users-number", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"nb_users": "9000"})
	})

	if useHTTPS {
		// Start an HTTP server to redirect traffic to HTTPS
		// go func() {
		// 	// TODO change the port to 80 or redirect it somehow via nginx
		// 	if err := http.ListenAndServe(":3373", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 		http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
		// 	})); err != nil {
		// 		log.Fatal("Failed to start HTTP redirect server:", err)
		// 	}
		// }()
		// Start the HTTPS server
		err = r.RunTLS(":443", "cert.pem", "key.pem")
		if err != nil {
			log.Fatal("Failed to start HTTPS server:", err)
		}
	} else {
		// Start the HTTP server
		//err = r.Run(":8080")
		//if err != nil {
		//	log.Fatal("Failed to start HTTP server:", err)
		//}
	}
}
