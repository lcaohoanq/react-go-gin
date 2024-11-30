package handlers

import (
	"github.com/lcaohoanq/react-go-gin/server/internal/constants"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lcaohoanq/react-go-gin/server/internal/models"
	"github.com/lcaohoanq/react-go-gin/server/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	input := new(RegisterInput)

	if err := c.ShouldBindJSON(input); err != nil {
		log.Printf("Error parsing input: %v", err)
		c.JSON(400, gin.H{
			"error": "Invalid input format",
		})
		return
	}

	match, _ := regexp.MatchString(constants.EmailRegex, input.Name)
	if !match {
		c.JSON(400, gin.H{
			"error":  "Invalid name format",
			"reason": "name must not contain any number, every first letter need capitalized, not contain any punctuation",
		})
		return
	}

	// Validate input
	if input.Name == "" || input.Email == "" || input.Password == "" {
		c.JSON(400, gin.H{
			"error": "All fields are required",
		})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(400, gin.H{
			"error": "Email already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Create user
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Avatar:   "https://ui-avatars.com/api/?name=" + input.Name,
		Role:     models.MEMBER,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		log.Printf("Error creating user: %v", result.Error)
		c.JSON(500, gin.H{
			"error": "Could not create user",
		})
		return
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(500, gin.H{
			"error": "Could not generate token",
		})
		return
	}

	c.JSON(201, gin.H{
		"token": t,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Login(c *gin.Context) {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(400, gin.H{
			"error":  "Invalid input",
			"reason": "Invalid input JSON format",
		})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"error":  "Wrong email or password",
			"reason": "User with this email does not exist",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(400, gin.H{
			"error":  "Wrong email or password",
			"reason": "Password is incorrect",
		})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Could not generate token",
		})
		return
	}

	c.JSON(200, gin.H{
		"token": t,
		"data": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"role":   user.Role,
			"avatar": user.Avatar,
		},
	})
}

func Logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Logout successful",
	})
}
