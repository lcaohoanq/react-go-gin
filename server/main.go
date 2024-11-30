package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lcaohoanq/react-go-gin/server/internal/routes"
	"github.com/lcaohoanq/react-go-gin/server/pkg/database"
	"log"
	"os"
)

func CORSMiddleware() gin.HandlerFunc {
	// Define allowed origins
	allowedOrigins := map[string]bool{
		"http://localhost:3000":  true,
		"localhost:3000":         true,
		"https://yourdomain.com": true,
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if the origin is in the allowed list
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()

	r := gin.Default()

	// Use the custom CORS middleware
	r.Use(CORSMiddleware())

	// Routes
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(r.Run(":" + port))
}
