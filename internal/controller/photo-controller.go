package controller

import (
	"fmt"
	"github.com/ZakyHermawan/maigram/common"
	"github.com/ZakyHermawan/maigram/entity"
	"github.com/ZakyHermawan/maigram/internal/repository"
	"github.com/ZakyHermawan/maigram/internal/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type PhotoController interface {
	CreatePhoto(ctx *gin.Context)
	GetAllPhoto(ctx *gin.Context)
	UpdatePhoto(ctx *gin.Context)
	DeletePhoto(ctx *gin.Context)
}

type photoController struct {
	photoService   service.PhotoService
	jwtService     service.JWTService
	userRepository repository.UserRepository
}

func NewPhotoController(photoService service.PhotoService, jwtService service.JWTService, userRepository repository.UserRepository) PhotoController {
	validate = validator.New()
	return &photoController{
		photoService:   photoService,
		jwtService:     jwtService,
		userRepository: userRepository,
	}
}

func (controller *photoController) CreatePhoto(ctx *gin.Context) {
	var photo entity.Photo
	if err := ctx.ShouldBindJSON(&photo); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to create photo",
		})
		return
	}
	if err := validate.Struct(&photo); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to create photo",
		})
		return
	}
	jwtToken, parseErr := common.ExtractBearerToken(ctx.GetHeader("Authorization"))
	if parseErr != nil {
		fmt.Println(parseErr.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": parseErr.Error(),
		})
		return
	}
	token, err := controller.jwtService.ValidateToken(jwtToken)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to validate token",
		})
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	user := entity.User{
		Email: email,
	}
	err = controller.userRepository.GetUserByEmail(&user)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to create photo",
		})
	}
	photo.UserID = user.ID
	err = controller.photoService.CreatePhoto(&photo)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to create photo",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"id":         photo.ID,
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoUrl,
		"user_id":    photo.UserID,
		"created_at": photo.CreatedAt,
	})
}

func (controller *photoController) GetAllPhoto(ctx *gin.Context) {
	photos, err := controller.photoService.GetAllPhoto()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to get all photo",
		})
		return
	}
	newArray := make([]map[string]interface{}, len(photos))
	for i, photo := range photos {
		newArray[i] = map[string]interface{}{}
		newArray[i]["id"] = photo.Model.ID
		newArray[i]["title"] = photo.Title
		newArray[i]["caption"] = photo.Caption
		newArray[i]["photo_url"] = photo.PhotoUrl
		newArray[i]["user_id"] = photo.UserID
		newArray[i]["created_at"] = photo.Model.CreatedAt
		newArray[i]["updated_at"] = photo.Model.UpdatedAt
		var user entity.User
		user.ID = photo.UserID
		err = controller.userRepository.GetUserByID(&user)
		if err != nil {
			fmt.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "fail to get all photo",
			})
			return
		}
		newArray[i]["User"] = map[string]string{
			"username": user.Username,
			"email":    user.Email,
		}
	}
	ctx.JSON(http.StatusOK, newArray)
}

func (controller *photoController) UpdatePhoto(ctx *gin.Context) {
	var tempPhoto = struct {
		Title    string `json:"title"`
		Caption  string `json:"caption"`
		PhotoUrl string `json:"photo_url"`
	}{}

	if err := ctx.ShouldBindJSON(&tempPhoto); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to edit photo",
		})
		return
	}

	id, parseErr := strconv.ParseUint(ctx.Param("photoId"), 0, 0)
	if parseErr != nil {
		fmt.Println(parseErr.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "fail to update user",
		})
	}

	photo := entity.Photo{
		Model: gorm.Model{
			ID: uint(id),
		},
		Title:    tempPhoto.Title,
		Caption:  tempPhoto.Caption,
		PhotoUrl: tempPhoto.PhotoUrl,
	}

	err := controller.photoService.UpdatePhoto(&photo)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to edit photo",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":         photo.ID,
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoUrl,
		"user_id":    photo.UserID,
		"updated_at": photo.UpdatedAt,
	})
}

func (controller *photoController) DeletePhoto(ctx *gin.Context) {
	id, parseErr := strconv.ParseUint(ctx.Param("photoId"), 0, 0)
	if parseErr != nil {
		fmt.Println(parseErr.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to delete photo",
		})
	}

	photo := entity.Photo{
		Model: gorm.Model{
			ID: uint(id),
		},
	}

	err := controller.photoService.DeletePhoto(&photo)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to delete photo",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}
