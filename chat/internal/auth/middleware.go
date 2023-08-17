package auth

import (
	"context"
	"fmt"
	"main/config"
	"net/http"
)

var cookieName = "auth-cookie"

type userClaims struct {
	id int64
}

type ContextKey string

var userClaimsKey = "user"

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(cookieName)

			if config.C.Env == "dev" {
				c = &http.Cookie{Value: "I am a cookie"}
			}

			if err != nil || c == nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}

			if config.C.Env == "dev" {
				ctx := context.WithValue(
					r.Context(), ContextKey(userClaimsKey), &userClaims{id: 11},
				)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUser(ctx context.Context) (*userClaims, error) {
	userClaims, ok := ctx.Value(ContextKey(userClaimsKey)).(*userClaims)
	if !ok {
		return nil, fmt.Errorf("Fail to extract user claims")
	}

	return userClaims, nil
}
