package controller

import (
	"fmt"
	"github.com/ZakyHermawan/maigram/common"
	"github.com/ZakyHermawan/maigram/entity"
	service2 "github.com/ZakyHermawan/maigram/internal/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type UserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	ChangeUsernameOrEmail(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type userController struct {
	userService service2.UserService
	jwtService  service2.JWTService
}

var validate *validator.Validate

func NewAuthController(authService service2.UserService, jwtService service2.JWTService) UserController {
	validate = validator.New()
	return &userController{
		userService: authService,
		jwtService:  jwtService,
	}
}

func (controller *userController) Register(ctx *gin.Context) {
	var user entity.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to create user",
		})
		return
	}
	if err := validate.Struct(&user); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to create user",
		})
		return
	}

	err := controller.userService.Register(&user)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to create user",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"age":      user.Age,
	})
}

func (controller *userController) Login(ctx *gin.Context) {
	tempUser := struct {
		Email    string
		Password string
	}{}
	if err := ctx.ShouldBindJSON(&tempUser); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "login failed",
		})
		return
	}
	var user entity.User
	user.Email = tempUser.Email
	user.Password = tempUser.Password
	isAuthenticated, err := controller.userService.Login(&user)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "login failed",
		})
		return
	}
	if !isAuthenticated {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid username or password",
		})
		return
	}
	token := controller.jwtService.GenerateToken(user.Email)
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (controller *userController) ChangeUsernameOrEmail(ctx *gin.Context) {
	tempUser := struct {
		Username string
		Email    string
	}{}
	if err := ctx.ShouldBindJSON(&tempUser); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to update user",
		})
		return
	}

	id, parseErr := strconv.ParseUint(ctx.Param("userId"), 0, 0)
	if parseErr != nil {
		fmt.Println(parseErr.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "fail to update user",
		})
	}

	user := entity.User{
		Model: gorm.Model{
			ID: uint(id),
		},
		Username: tempUser.Username,
		Email:    tempUser.Email,
	}

	err := controller.userService.UpdateUsernameOrEmail(&user)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to update user",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"age":        user.Age,
		"updated_at": user.UpdatedAt,
	})
}

func (controller *userController) DeleteUser(ctx *gin.Context) {
	jwtToken, parseErr := common.ExtractBearerToken(ctx.GetHeader("Authorization"))
	if parseErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": parseErr.Error(),
		})
		return
	}
	token, err := controller.jwtService.ValidateToken(jwtToken)
	if err != nil {
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
	err = controller.userService.DeleteUser(&user)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to delete user",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}
