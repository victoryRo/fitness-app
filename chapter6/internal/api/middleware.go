package api

import (
	"mime"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
)

// JSONMiddleware nos ayudará a manejar sólo JSON
// dentro y fuera
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if strings.TrimSpace(contentType) == "" {
			var parseError error
			contentType, _, parseError = mime.ParseMediaType(contentType)
			if parseError != nil {
				JSONError(w, http.StatusBadRequest, "Bad or not Content-Type header found")
				return
			}
		}

		if contentType != "application/json" {
			JSONError(w, http.StatusUnsupportedMediaType, "Content-Type not application/json")
			return
		}

		// Dígale al cliente que también estamos hablando de JSON
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(origins []string) func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{
			"Accept",
			"Origin",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Methods",
		}),

		handlers.AllowedOrigins(origins),
		handlers.AllowedMethods([]string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPost,
			http.MethodPatch,
			http.MethodDelete,
		}),
	)
}
