package middleware

import (
	"log"
	"messenger/internal/service"
	"net/http"
	"strings"
)

func Auth(s *service.Service, next http.Handler) http.Handler {
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

func MakeAuthMiddleware(s *service.Service, tokenLifeTimeMinutes int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			endpoint := r.URL.String()

			if endpoint != "/auth" && endpoint[:4] != "/ws?" {
				token, ok := strings.CutPrefix(r.Header.Get("Authorization"),
					"Bearer ",
				)
				if !ok {
					log.Println(endpoint)
					http.Error(w, "invalid authorization header", http.StatusUnauthorized)
					return
				}
				res, err := s.CheckToken(r.Context(), token, tokenLifeTimeMinutes)
				if !res || err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"error": "unauthorized"}`))
					return // Прерываем выполнение
				}
			} else {
				//log.Println("это ауф запрос")
			}

			next.ServeHTTP(w, r)
		})
	}
}
