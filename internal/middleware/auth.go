package middleware

import (
	"log"
	"messenger/internal/service"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func MatchesPattern(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Проверяем путь: /users/{id}/pic или /channels/{id}/pic
	pathPattern := regexp.MustCompile(`^/(users|channels)/([^/]+)/pic$`)
	if !pathPattern.MatchString(parsedURL.Path) {
		return false
	}

	// Проверяем наличие token
	token := parsedURL.Query().Get("token")
	if token == "" {
		return false
	}

	// Проверяем, что нет лишних query-параметров
	if len(parsedURL.Query()) != 1 {
		return false
	}

	return true
}

// ExtractParams извлекает тип (users/channels), id и token из URL
func ExtractParams(rawURL string) (resourceType, id, token string, ok bool) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", "", false
	}

	re := regexp.MustCompile(`^/(users|channels)/([^/]+)/pic$`)
	matches := re.FindStringSubmatch(parsedURL.Path)
	if matches == nil {
		return "", "", "", false
	}

	resourceType = matches[1] // "users" или "channels"
	id = matches[2]           // user_id или channel_id

	token = parsedURL.Query().Get("token")
	if token == "" {
		return "", "", "", false
	}

	if len(parsedURL.Query()) != 1 {
		return "", "", "", false
	}

	return resourceType, id, token, true
}

// ExtractParamsDecoded извлекает и декодирует параметры
func ExtractParamsDecoded(rawURL string) (resourceType, id, token string, ok bool) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", "", false
	}

	re := regexp.MustCompile(`^/(users|channels)/([^/]+)/pic$`)
	matches := re.FindStringSubmatch(parsedURL.Path)
	if matches == nil {
		return "", "", "", false
	}

	resourceType = matches[1]

	id, err = url.PathUnescape(matches[2])
	if err != nil {
		return "", "", "", false
	}

	token = parsedURL.Query().Get("token")
	if token == "" {
		return "", "", "", false
	}

	if len(parsedURL.Query()) != 1 {
		return "", "", "", false
	}

	return resourceType, id, token, true
}

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

			if endpoint == "/auth" {
				next.ServeHTTP(w, r)
			} else if endpoint[:4] == "/ws?" || MatchesPattern(endpoint) {
				token := strings.TrimSpace(r.URL.Query().Get("token"))

				res, err := s.CheckToken(r.Context(), token, tokenLifeTimeMinutes)
				if !res || err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"error": "unauthorized"}`))
					return // Прерываем выполнение
				}
				next.ServeHTTP(w, r)

			} else {
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
				next.ServeHTTP(w, r)
			}
		})
	}
}
