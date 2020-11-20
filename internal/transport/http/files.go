package http

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

func (h *handler) handleUploadFile() http.HandlerFunc {
	type Response struct {
		withStatusCreated
		FileName string `json:"fileName"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(contextKeyUserUUID)
		if userID == nil {
			respond(w, r, unauthorized("unauthorized request"))
			return
		}

		res, err := h.fileSvc.UploadFile(r.Context(), r.Body)
		if err != nil {
			respond(w, r, internalServerError(errors.Wrap(err, "error on upload file")))
			return
		}

		rsp := Response{FileName: res.FileName}

		respond(w, r, rsp)
	}
}

func (h *handler) handleDownloadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := mux.Vars(r)["fileName"]
		res, err := h.fileSvc.DownloadFile(r.Context(), fileName)
		if err != nil {
			respond(w, r, internalServerError(errors.Wrap(err, "error on download file")))
			return
		}

		lm, err := h.fileSvc.GetFileLastModified(r.Context(), fileName)
		if err != nil {
			respond(w, r, internalServerError(errors.Wrap(err, "error on get file last modified")))
			return
		}

		http.ServeContent(w, r, fileName, *lm, res)
	}
}
