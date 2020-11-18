package http

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func cors() mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{
			http.MethodOptions,
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		}),
		handlers.AllowedHeaders([]string{
			"Authorization",
			"Content-Type",
			"Content-Language",
			"Accept",
			"Accept-Language",
			"Origin",
		}),
		handlers.AllowCredentials(),
	)
}
