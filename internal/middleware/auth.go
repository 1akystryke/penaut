package middleware

import (
	"log"
	"net/http"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AUTHORIZATION")
		endpoint := r.URL.String()

		token, ok := strings.CutPrefix(r.Header.Get("Authorization"),
			"Bearer ",
		)

		if !ok {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}
		log.Printf("%s %s %s", r.Method, endpoint, token)
		next.ServeHTTP(w, r)
	})
}
