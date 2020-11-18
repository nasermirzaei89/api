package http

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

type recoverMW struct {
	next http.Handler
}

func (mw *recoverMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if v := recover(); v != nil {
			err := errors.Errorf("panic recovered: %+v", v)
			rsp := internalServerError(err)
			respond(w, r, rsp)
		}
	}()
	mw.next.ServeHTTP(w, r)
}

func recoverPanic() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return &recoverMW{
			next: next,
		}
	}
}
