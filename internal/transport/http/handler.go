package http

import (
	gogzip "compress/gzip"
	"github.com/gorilla/mux"
	"github.com/nasermirzaei89/api/internal/services/file"
	"github.com/nasermirzaei89/api/internal/services/post"
	"github.com/nasermirzaei89/api/internal/services/user"
	"net/http"
)

type handler struct {
	router  *mux.Router
	userSvc user.Service
	postSvc post.Service
	fileSvc file.Service
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func NewHandler(l loggerInterface, userSvc user.Service, postSvc post.Service, fileSvc file.Service) http.Handler {
	r := mux.NewRouter()

	r.Use(cors())
	r.Use(gzip(gogzip.BestSpeed))
	r.Use(logger(l))
	r.Use(recoverPanic())
	r.Use(authenticate(userSvc))

	h := handler{
		router:  r,
		userSvc: userSvc,
		postSvc: postSvc,
		fileSvc: fileSvc,
	}

	h.router.Methods(http.MethodPost, http.MethodGet).Path("/graphql").HandlerFunc(h.handleGraphQL())
	h.router.Methods(http.MethodPost).Path("/files").HandlerFunc(h.handleUploadFile())
	h.router.Methods(http.MethodGet).Path("/files/{fileName}").HandlerFunc(h.handleDownloadFile())

	return &h
}
