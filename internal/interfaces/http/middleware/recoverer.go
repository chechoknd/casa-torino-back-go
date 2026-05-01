package middleware

import (
	"log"
	"net/http"

	"github.com/casatorino/backend/internal/interfaces/http/responses"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic recovered: %v", recovered)
				responses.WriteError(w, nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
