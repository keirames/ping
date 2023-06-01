package jwt

import (
	"context"
	"fmt"
	"main/tools"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaim struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

func GetClaims(ctx *gin.Context) *JwtCustomClaim {
	value, ok := ctx.Get("claims")
	if !ok {
		return nil
	}

	claims, ok := value.(*JwtCustomClaim)
	if !ok {
		return nil
	}

	return claims
}

func GenerateJwt(ctx context.Context, userID int64) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	fmt.Println([]byte(tools.Config.JWTSecret))
	token, err := t.SignedString([]byte(tools.Config.JWTSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateJwt(ctx context.Context, token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &JwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}

		return []byte(tools.Config.JWTSecret), nil
	})
}
