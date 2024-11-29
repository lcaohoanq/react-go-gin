package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcaohoanq/react-go-gin/server/internal/models"
	"github.com/lcaohoanq/react-go-gin/server/pkg/database"
)

func GetTodos(c *gin.Context) {
	userID := c.GetFloat64("user_id")

	var todos []models.Todo
	result := database.DB.Where("user_id = ?", uint(userID)).Find(&todos)
	if result.Error != nil {
		log.Printf("Error fetching todos: %v", result.Error)
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch todos",
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, todos)
}

func CreateTodo(c *gin.Context) {
	userID := c.GetFloat64("user_id")

	var input struct {
		Body string `json:"body"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	if input.Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Todo body cannot be empty",
		})
		return
	}

	todo := models.Todo{
		Body:   input.Body,
		UserID: uint(userID),
	}

	result := database.DB.Create(&todo)
	if result.Error != nil {
		log.Printf("Error creating todo: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create todo",
		})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func UpdateTodo(c *gin.Context) {
	userID := c.GetFloat64("user_id")
	id := c.Param("id")

	var todo models.Todo
	if err := database.DB.Where("id = ? AND user_id = ?", id, uint(userID)).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Todo not found or unauthorized",
		})
		return
	}

	todo.Completed = true
	if err := database.DB.Save(&todo).Error; err != nil {
		log.Printf("Error updating todo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update todo",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context) {
	userID := c.GetFloat64("user_id")
	id := c.Param("id")

	var todo models.Todo
	if err := database.DB.Where("id = ? AND user_id = ?", id, uint(userID)).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Todo not found or unauthorized",
		})
		return
	}

	if err := database.DB.Delete(&todo).Error; err != nil {
		log.Printf("Error deleting todo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete todo",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
