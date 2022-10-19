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

type CommentController interface {
	CreateComment(ctx *gin.Context)
	GetAllComment(ctx *gin.Context)
	UpdateComment(ctx *gin.Context)
	DeleteComment(ctx *gin.Context)
}

type commentController struct {
	commentService  service.CommentService
	jwtService      service.JWTService
	userRepository  repository.UserRepository
	photoRepository repository.PhotoRepository
}

func NewCommentController(
	commentService service.CommentService,
	jwtService service.JWTService,
	userRepository repository.UserRepository,
	photoRepository repository.PhotoRepository,
) CommentController {
	validate = validator.New()
	return &commentController{
		commentService:  commentService,
		jwtService:      jwtService,
		userRepository:  userRepository,
		photoRepository: photoRepository,
	}
}

func (controller *commentController) CreateComment(ctx *gin.Context) {
	var comment entity.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to create comment",
		})
		return
	}
	if err := validate.Struct(&comment); err != nil {
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
	comment.UserID = user.ID

	photo := entity.Photo{
		Model: gorm.Model{
			ID: comment.PhotoID,
		},
	}
	err = controller.photoRepository.GetPhotoById(&photo)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to create comment",
		})
		return
	}

	err = controller.commentService.CreateComment(&comment)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to create comment",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"id":         comment.ID,
		"message":    comment.Message,
		"photo_id":   comment.PhotoID,
		"user_id":    comment.UserID,
		"created_at": comment.CreatedAt,
	})
}

func (controller *commentController) GetAllComment(ctx *gin.Context) {
	comments, err := controller.commentService.GetAllComment()
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to get all comment",
		})
		return
	}
	newArray := make([]map[string]interface{}, len(comments))
	for i, comment := range comments {
		newArray[i] = map[string]interface{}{}
		newArray[i]["id"] = comment.Model.ID
		newArray[i]["message"] = comment.Message
		newArray[i]["photo_id"] = comment.PhotoID
		newArray[i]["user_id"] = comment.UserID
		newArray[i]["created_at"] = comment.Model.CreatedAt
		newArray[i]["updated_at"] = comment.Model.UpdatedAt
		var user entity.User
		user.ID = comment.UserID
		err = controller.userRepository.GetUserByID(&user)
		if err != nil {
			fmt.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "fail to get all comment",
			})
			return
		}
		newArray[i]["User"] = map[string]interface{}{
			"id":       user.Model.ID,
			"username": user.Username,
			"email":    user.Email,
		}
		var photo entity.Photo
		photo.ID = comment.PhotoID
		err = controller.photoRepository.GetPhotoById(&photo)
		if err != nil {
			fmt.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "fail to get all comment",
			})
			return
		}
		newArray[i]["Photo"] = map[string]interface{}{
			"id":        photo.ID,
			"title":     photo.Title,
			"caption":   photo.Caption,
			"photo_url": photo.PhotoUrl,
			"user_id":   photo.UserID,
		}
	}
	ctx.JSON(http.StatusOK, newArray)
}

func (controller *commentController) UpdateComment(ctx *gin.Context) {
	var tempComment = struct {
		Message string `json:"message"`
	}{}
	if err := ctx.ShouldBindJSON(&tempComment); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to update comment",
		})
		return
	}
	if err := validate.Struct(&tempComment); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to update comment",
		})
		return
	}

	id, parseParamErr := strconv.ParseUint(ctx.Param("commentId"), 0, 0)
	if parseParamErr != nil {
		fmt.Println(parseParamErr.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "fail to update comment",
		})
	}

	comment := entity.Comment{
		Model: gorm.Model{
			ID: uint(id),
		},
		Message: tempComment.Message,
	}
	err := controller.commentService.UpdateComment(&comment)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to update comment",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         comment.ID,
		"message":    comment.Message,
		"photo_id":   comment.PhotoID,
		"user_id":    comment.UserID,
		"updated_at": comment.UpdatedAt,
	})
}

func (controller *commentController) DeleteComment(ctx *gin.Context) {
	id, parseParamErr := strconv.ParseUint(ctx.Param("commentId"), 0, 0)
	if parseParamErr != nil {
		fmt.Println(parseParamErr.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to delete comment",
		})
	}
	comment := entity.Comment{
		Model: gorm.Model{
			ID: uint(id),
		},
	}
	err := controller.commentService.DeleteComment(&comment)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to delete comment",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
