package auth

import (
	"fmt"
	"net/http"
)

var cookieName = "auth-cookie"

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(cookieName)

			fmt.Println(c)
			// Allow unauthenticated users in
			if err != nil || c == nil {
				fmt.Println("inside here", err, c)
				next.ServeHTTP(w, r)
				return
			}
			fmt.Println("got cookie", c)

			next.ServeHTTP(w, r)
		})
	}
}
