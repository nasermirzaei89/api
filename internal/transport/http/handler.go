package http

import (
	"github.com/gorilla/mux"
	"github.com/nasermirzaei89/api/internal/services/post"
	"github.com/nasermirzaei89/api/internal/services/user"
	"net/http"
)

type handler struct {
	router  *mux.Router
	userSvc user.Service
	postSvc post.Service
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func NewHandler(l loggerInterface, userSvc user.Service, postSvc post.Service) http.Handler {
	r := mux.NewRouter()

	r.Use(cors())
	r.Use(logger(l))
	r.Use(recoverPanic())
	r.Use(authenticate(userSvc))

	h := handler{
		router:  r,
		userSvc: userSvc,
		postSvc: postSvc,
	}

	h.router.Methods(http.MethodPost).Path("/graphql").HandlerFunc(h.handleGraphQL())

	return &h
}
