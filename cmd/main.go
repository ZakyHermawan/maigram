package main

import (
	"fmt"
	"github.com/ZakyHermawan/maigram/common"
	"github.com/ZakyHermawan/maigram/entity"
	"github.com/ZakyHermawan/maigram/internal/middleware"
	"github.com/ZakyHermawan/maigram/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Photo{})
	db.AutoMigrate(&entity.Comment{})
	db.AutoMigrate(&entity.SocialMedia{})
}

func main() {
	db := common.Init()
	defer func() { _ = db.Close() }()

	Migrate(db)

	server := gin.Default()

	router.InitRouter()
	router.UserRegister(server.Group("/users"))
	router.PhotoRegister(server.Group("/photos", middleware.AuthorizeJWT()))
	router.CommentRegister(server.Group("/comments", middleware.AuthorizeJWT()))
	router.SocialMediaRegister(server.Group("/socialmedias", middleware.AuthorizeJWT()))

	if err := server.Run("localhost:8080"); err != nil {
		fmt.Println(err.Error())
		return
	}
}
