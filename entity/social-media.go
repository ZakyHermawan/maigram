package entity

import "github.com/jinzhu/gorm"

type SocialMedia struct {
	gorm.Model
	Name           string `json:"name" binding:"required" gorm:"type:varchar(255)"`
	SocialMediaURL string `json:"social_media_url" binding:"required,url" gorm:"type:varchar(255)"`
	UserID         uint   `json:"user_id"`
}
