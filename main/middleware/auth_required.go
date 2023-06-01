package middleware

import (
	"context"
	"fmt"
	"main/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := ctx.Cookie("auth-cookie")

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		newToken, _ := jwt.GenerateJwt(context.Background(), 1)

		validate, err := jwt.ValidateJwt(context.Background(), newToken)
		if err != nil {
			fmt.Println(err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		customClaims, ok := validate.Claims.(*jwt.JwtCustomClaim)
		if !ok {
			fmt.Println("fail to extract jwt custom claim")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("claims", customClaims)

		ctx.Next()
	}
}
