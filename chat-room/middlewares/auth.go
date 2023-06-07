package middlewares

import (
	"chatroom/logger"
	"context"
	"fmt"
	"net/http"
)

type ctxKey string

func GetUserID(ctx context.Context) int64 {
	return 61843283797778460

	userID, ok := ctx.Value(ctxKey("userID")).(int64)
	if !ok {
		logger.L.Error().Msg("Cannot extract token value in context")
		return 0
	}

	return userID
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
		ctx := context.WithValue(r.Context(), ctxKey("userID"), "1")
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
