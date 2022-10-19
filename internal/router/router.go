package router

import (
	"github.com/ZakyHermawan/maigram/internal/controller"
	"github.com/ZakyHermawan/maigram/internal/middleware"
	"github.com/ZakyHermawan/maigram/internal/repository"
	"github.com/ZakyHermawan/maigram/internal/service"
	"github.com/gin-gonic/gin"
)

var userRepository repository.UserRepository
var photoRepository repository.PhotoRepository
var commentRepository repository.CommentRepository
var socialMediaRepository repository.SocialMediaRepository

var userService service.UserService
var photoService service.PhotoService
var commentService service.CommentService
var socialMediaService service.SocialMediaService
var jwtService service.JWTService

var userController controller.UserController
var photoController controller.PhotoController
var commentController controller.CommentController
var socialMediaController controller.SocialMediaController

func InitRouter() {
	userRepository = repository.NewUserRepository()
	photoRepository = repository.NewPhotoRepository()
	commentRepository = repository.NewCommentRepository()
	socialMediaRepository = repository.NewSocialMediaRepository()

	userService = service.NewUserService(userRepository)
	photoService = service.NewPhotoService(photoRepository)
	commentService = service.NewCommentService(commentRepository)
	socialMediaService = service.NewSocialMediaService(socialMediaRepository)
	jwtService = service.NewJWTService()

	userController = controller.NewAuthController(userService, jwtService)
	photoController = controller.NewPhotoController(photoService, jwtService, userRepository)
	commentController = controller.NewCommentController(commentService, jwtService, userRepository, photoRepository)
	socialMediaController = controller.NewSocialMediaController(socialMediaService, jwtService, userRepository)
}

func UserRegister(router *gin.RouterGroup) {
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.Use(middleware.AuthorizeJWT())
	router.PUT("/:userId", userController.ChangeUsernameOrEmail)
	router.DELETE("", userController.DeleteUser)
}

func PhotoRegister(router *gin.RouterGroup) {
	router.POST("", photoController.CreatePhoto)
	router.GET("", photoController.GetAllPhoto)
	router.PUT("/:photoId", photoController.UpdatePhoto)
	router.DELETE("/:photoId", photoController.DeletePhoto)
}

func CommentRegister(router *gin.RouterGroup) {
	router.POST("", commentController.CreateComment)
	router.GET("", commentController.GetAllComment)
	router.PUT("/:commentId", commentController.UpdateComment)
	router.DELETE("/:commentId", commentController.DeleteComment)
}

func SocialMediaRegister(router *gin.RouterGroup) {
	router.POST("", socialMediaController.CreateSocialMedia)
	router.GET("", socialMediaController.GetAllSocialMedia)
	router.PUT("/:socialMediaId", socialMediaController.UpdateSocialMedia)
	router.DELETE("/:socialMediaId", socialMediaController.DeleteSocialMedia)
}
