package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lcaohoanq/react-go-gin/server/internal/handlers"
	"github.com/lcaohoanq/react-go-gin/server/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Public routes
	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)

	// Protected routes
	todos := api.Group("/todos")
	todos.Use(middleware.Protected())
	{
		todos.GET("/", handlers.GetTodos)
		todos.POST("/", handlers.CreateTodo)
		todos.PATCH("/:id", handlers.UpdateTodo)
		todos.DELETE("/:id", handlers.DeleteTodo)
	}

	api.GET("/profile", middleware.Protected(), handlers.GetUserProfile)
}
