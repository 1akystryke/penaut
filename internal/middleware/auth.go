package middleware

import (
	"log"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AUTHORIZATION")
		endpoint := r.URL.String()
		header := r.Header["Authorization"]
		log.Printf("%s %s %s", r.Method, endpoint, header)
		next.ServeHTTP(w, r)
	})
}
