package jwt

import (
	"context"
	"fmt"
	"time"

	"chatroom/config"

	"github.com/golang-jwt/jwt/v4"
)

type JWTCustomClaim struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateJwt(ctx context.Context, userID string) (*string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTCustomClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err := t.SignedString([]byte(config.C.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func ValidateJwt(ctx context.Context, token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &JWTCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}

		return []byte(config.C.JWTSecret), nil
	})
}
