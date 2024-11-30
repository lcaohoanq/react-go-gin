package models

import (
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Body      string `json:"body"`
	Completed bool   `json:"completed" gorm:"default:false"`
	UserID    uint   `json:"userId"`
	User      User   `json:"-" gorm:"foreignKey:UserID"` // The `-` will exclude this field from JSON
}
