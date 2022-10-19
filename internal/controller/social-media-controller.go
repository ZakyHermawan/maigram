package controller

import (
	"fmt"
	"github.com/ZakyHermawan/maigram/common"
	"github.com/ZakyHermawan/maigram/entity"
	"github.com/ZakyHermawan/maigram/internal/repository"
	"github.com/ZakyHermawan/maigram/internal/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type SocialMediaController interface {
	CreateSocialMedia(ctx *gin.Context)
	GetAllSocialMedia(ctx *gin.Context)
	UpdateSocialMedia(ctx *gin.Context)
	DeleteSocialMedia(ctx *gin.Context)
}

type socialMediaController struct {
	socialMediaService service.SocialMediaService
	jwtService         service.JWTService
	userRepository     repository.UserRepository
}

func NewSocialMediaController(
	socialMediaService service.SocialMediaService,
	jwtService service.JWTService,
	userRepository repository.UserRepository,
) SocialMediaController {
	return &socialMediaController{
		socialMediaService: socialMediaService,
		jwtService:         jwtService,
		userRepository:     userRepository,
	}
}

func (controller *socialMediaController) CreateSocialMedia(ctx *gin.Context) {
	var socialMedia entity.SocialMedia
	if err := ctx.ShouldBindJSON(&socialMedia); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to create comment",
		})
		return
	}

	if err := validate.Struct(&socialMedia); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to create comment",
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
			"message": "fail to create comment",
		})
		return
	}

	socialMedia.UserID = user.ID
	err = controller.socialMediaService.CreateSocialMedia(&socialMedia)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to create comment",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":               socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id":          socialMedia.UserID,
		"created_at":       socialMedia.CreatedAt,
	})
}

func (controller *socialMediaController) GetAllSocialMedia(ctx *gin.Context) {
	socialMedias, err := controller.socialMediaService.GetAllSocialMedia()
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to get all social media",
		})
		return
	}

	newArray := make([]map[string]interface{}, len(socialMedias))
	for i, socialMedia := range socialMedias {
		newArray[i] = map[string]interface{}{
			"id":               socialMedia.ID,
			"name":             socialMedia.Name,
			"social_media_url": socialMedia.SocialMediaURL,
			"user_id":          socialMedia.UserID,
			"created_at":       socialMedia.CreatedAt,
			"updated_at":       socialMedia.UpdatedAt,
		}
		var user entity.User
		user.ID = socialMedia.UserID
		err = controller.userRepository.GetUserByID(&user)
		if err != nil {
			fmt.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "fail to get all social media",
			})
			return
		}
		newArray[i]["user"] = map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			// possible changes on this response
			//"email":    user.Email,
		}
	}
	ctx.JSON(http.StatusOK, newArray)
}

func (controller *socialMediaController) UpdateSocialMedia(ctx *gin.Context) {
	var tempSocialMedia = struct {
		Name           string `json:"name"`
		SocialMediaURL string `json:"social_media_url"`
	}{}
	if err := ctx.ShouldBindJSON(&tempSocialMedia); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to update social media",
		})
		return
	}

	if err := validate.Struct(&tempSocialMedia); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to update social media",
		})
		return
	}

	id, parseParamErr := strconv.ParseUint(ctx.Param("socialMediaId"), 0, 0)
	if parseParamErr != nil {
		fmt.Println(parseParamErr.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "fail to update comment",
		})
		return
	}

	socialMedia := entity.SocialMedia{
		Model: gorm.Model{
			ID: uint(id),
		},
		Name:           tempSocialMedia.Name,
		SocialMediaURL: tempSocialMedia.SocialMediaURL,
	}

	err := controller.socialMediaService.UpdateSocialMedia(&socialMedia)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to update social media",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":               socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id":          socialMedia.UserID,
		"updated_at":       socialMedia.UpdatedAt,
	})
}

func (controller *socialMediaController) DeleteSocialMedia(ctx *gin.Context) {
	id, parseParamErr := strconv.ParseUint(ctx.Param("socialMediaId"), 0, 0)
	if parseParamErr != nil {
		fmt.Println(parseParamErr.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "fail to delete comment",
		})
		return
	}

	socialMedia := entity.SocialMedia{
		Model: gorm.Model{
			ID: uint(id),
		},
	}
	err := controller.socialMediaService.DeleteSocialMedia(&socialMedia)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to delete social media",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}
