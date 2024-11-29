package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lcaohoanq/react-go-gin/server/internal/models"
	"github.com/lcaohoanq/react-go-gin/server/pkg/database"
)

type UserStats struct {
	TotalTodos     int64         `json:"totalTodos"`
	CompletedTodos int64         `json:"completedTodos"`
	PendingTodos   int64         `json:"pendingTodos"`
	CompletionRate float64       `json:"completionRate"`
	RecentTodos    []models.Todo `json:"recentTodos"`
}

func GetUserProfile(c *gin.Context) {
	userId := c.GetFloat64("user_id")

	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "User not found",
		})
		return
	}

	// Get user statistics
	var stats UserStats

	// Total todos
	database.DB.Model(&models.Todo{}).Where("user_id = ?", userId).Count(&stats.TotalTodos)

	// Completed todos
	database.DB.Model(&models.Todo{}).Where("user_id = ? AND completed = ?", userId, true).Count(&stats.CompletedTodos)

	// Pending todos
	stats.PendingTodos = stats.TotalTodos - stats.CompletedTodos

	// Completion rate
	if stats.TotalTodos > 0 {
		stats.CompletionRate = float64(stats.CompletedTodos) / float64(stats.TotalTodos) * 100
	}

	// Recent todos (last 5)
	database.DB.Where("user_id = ?", userId).
		Order("created_at desc").
		Limit(5).
		Find(&stats.RecentTodos)

	c.JSON(200, gin.H{
		"user":  user,
		"stats": stats,
	})
}
