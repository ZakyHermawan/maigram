package entity

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username        string `json:"username" binding:"required" gorm:"type:varchar(255);UNIQUE"`
	Password        string `json:"password" binding:"required,min=6" gorm:"type:varchar(255)"`
	Email           string `json:"email" binding:"required,email" gorm:"type:varchar(255);UNIQUE"`
	Age             int    `json:"age" binding:"required,gt=8"`
	ProfileImageUrl string `json:"profile_image_url" gorm:"type:varchar(255);NOT NULL"`
}
