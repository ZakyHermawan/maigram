package middleware

import (
	"fmt"
	"github.com/ZakyHermawan/maigram/common"
	"github.com/ZakyHermawan/maigram/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		jwtToken, parseErr := common.ExtractBearerToken(ctx.GetHeader("Authorization"))
		if parseErr != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": parseErr.Error(),
			})
			return
		}

		token, err := service.NewJWTService().ValidateToken(jwtToken)
		if err != nil {
			fmt.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Fail to validate token",
			})
		}
		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
