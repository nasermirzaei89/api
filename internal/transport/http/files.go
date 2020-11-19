package http

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func (h *handler) handleUploadFile() http.HandlerFunc {
	type Response struct {
		withStatusCreated
		FileName string `json:"fileName"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
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

		http.ServeContent(w, r, fileName, time.Time{}, res)
	}
}
