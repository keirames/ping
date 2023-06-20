package middlewares

import (
	"chatroom/jwt"
	"chatroom/logger"
	"context"
	"fmt"
	"net/http"
)

type ctxKey string

func GetUserID(ctx context.Context) int64 {
	claims, ok := ctx.Value(ctxKey("userID")).(*jwt.JWTCustomClaim)
	if !ok {
		// TODO: don't know what to do if this occur
		logger.L.Error().Msg("Cannot extract token value in context")
		return 0
	}

	return claims.UserID
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hc, err := r.Cookie("x-token")
		if err != nil {
			logger.L.Error().Err(err).Msg("Cannot get cookie")

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Println(hc.Value)
		t, err := jwt.ValidateJwt(context.Background(), hc.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKey("userID"), t.Claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func VerifyTokenValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())

		fmt.Println("userID", userID)

		next.ServeHTTP(w, r)
	})
}
