package auth

import (
	"context"
	"fmt"
	"main/config"
	"net/http"
	"strconv"
)

var cookieName = "auth-cookie"

type userClaims struct {
	ID int64
}

type ContextKey string

var userClaimsKey = "user"

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.C.Env == "dev" {
				id := r.Header.Get("x-dev-user-id")
				userID, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					http.Error(w, "Invalid cookie", http.StatusForbidden)
					return
				}

				ctx := context.WithValue(
					r.Context(), ContextKey(userClaimsKey), &userClaims{ID: userID},
				)
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
				return
			}

			panic("not implemented")
			// if err != nil || c == nil {
			// 	http.Error(w, "Invalid cookie", http.StatusForbidden)
			// 	return
			// }
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
