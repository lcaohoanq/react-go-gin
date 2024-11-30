package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lcaohoanq/react-go-gin/server/internal/handlers"
	"github.com/lcaohoanq/react-go-gin/server/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Public routes

	auth := api.Group("/auth")
	auth.Use()
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		//auth.POST("/logout", handlers.Logout)
	}

	// Protected routes
	todos := api.Group("/todos")
	todos.Use(middleware.Protected())
	{
		todos.GET("/", handlers.GetTodos)
		todos.POST("/", handlers.CreateTodo)
		todos.PATCH("/:id", handlers.UpdateTodo)
		todos.DELETE("/:id", handlers.DeleteTodo)
	}

	users := api.Group("/users")
	users.Use(middleware.Protected())
	{
		users.GET("/profile", handlers.GetUserProfile)
	}
}
