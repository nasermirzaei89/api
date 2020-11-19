package http

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func gzip(level int) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return handlers.CompressHandlerLevel(h, level)
	}
}
