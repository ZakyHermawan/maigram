package entity

import "github.com/jinzhu/gorm"

type Comment struct {
	gorm.Model
	Message string `json:"message" binding:"required" gorm:"type:varchar(255)"`
	UserID  uint   `json:"user_id"`
	PhotoID uint   `json:"photo_id"`
}
