package http

import (
	"github.com/gorilla/mux"
	"github.com/nasermirzaei89/api/internal/services/file"
	"github.com/nasermirzaei89/api/internal/services/post"
	"github.com/nasermirzaei89/api/internal/services/user"
	"net/http"
)

type handler struct {
	router                  *mux.Router
	userSvc                 user.Service
	postSvc                 post.Service
	fileSvc                 file.Service
	enableGraphQLPretty     bool
	enableGraphQLPlayground bool
	enableGraphiQL          bool
	gzipLevel               int
	logger                  loggerInterface
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func NewHandler(l loggerInterface, userSvc user.Service, postSvc post.Service, fileSvc file.Service, options ...Option) http.Handler {
	h := handler{
		router:  mux.NewRouter(),
		userSvc: userSvc,
		postSvc: postSvc,
		fileSvc: fileSvc,
		logger:  l,
	}

	for i := range options {
		options[i](&h)
	}

	h.router.Use(cors())
	h.router.Use(gzip(h.gzipLevel))
	h.router.Use(logger(h.logger))
	h.router.Use(recoverPanic())
	h.router.Use(authenticate(h.userSvc))

	h.router.Path("/graphql").Handler(h.handleGraphQL(h.enableGraphQLPretty, h.enableGraphiQL, h.enableGraphQLPlayground))
	h.router.Methods(http.MethodPost).Path("/files").HandlerFunc(h.handleUploadFile())
	h.router.Methods(http.MethodGet).Path("/files/{fileName}").HandlerFunc(h.handleDownloadFile())

	return &h
}

type Option func(h *handler)

func SetGZipLevel(v int) Option {
	return func(h *handler) {
		h.gzipLevel = v
	}
}

func SetGraphQLPretty(v bool) Option {
	return func(h *handler) {
		h.enableGraphQLPretty = v
	}
}

func SetGraphiQL(v bool) Option {
	return func(h *handler) {
		h.enableGraphiQL = v
	}
}

func SetGraphQLPlayground(v bool) Option {
	return func(h *handler) {
		h.enableGraphQLPlayground = v
	}
}
