package http

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nasermirzaei89/api/internal/services/user"
	"net/http"
	"strings"
)

type contextKey string

const contextKeyUserUUID contextKey = "userUUID"

type authMW struct {
	next    http.Handler
	userSvc user.Service
}

func authenticate(userSvc user.Service) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return &authMW{
			next:    next,
			userSvc: userSvc,
		}
	}
}

func (mw *authMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		mw.next.ServeHTTP(w, r)
		return
	}

	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		respond(w, r, unauthorized("unsupported authorization header"))
		return
	}

	tokenString := authHeader[7:]

	usr, err := mw.userSvc.GetUserByTokenString(r.Context(), tokenString)
	if err != nil {
		respond(w, r, unauthorized("invalid authorization header"))
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), contextKeyUserUUID, usr.UUID))

	mw.next.ServeHTTP(w, r)
}
