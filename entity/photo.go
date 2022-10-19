package entity

import "github.com/jinzhu/gorm"

type Photo struct {
	gorm.Model
	Title    string `json:"title" binding:"required" gorm:"type:varchar(255)"`
	Caption  string `json:"caption" gorm:"type:varchar(255)"`
	PhotoUrl string `json:"photo_url" binding:"required,url" gorm:"type:varchar(255)"`
	UserID   uint   `json:"user_id"`
}
