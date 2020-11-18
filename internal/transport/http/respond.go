package http

import (
	"net/http"
)

func respond(w http.ResponseWriter, _ *http.Request, rsp interface{}) {
	if rsp == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if withHeader, ok := rsp.(interface{ Header() http.Header }); ok {
		for k, vv := range withHeader.Header() {
			for _, v := range vv {
				w.Header().Set(k, v)
			}
		}
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	if withStatusCode, ok := rsp.(interface{ StatusCode() int }); ok {
		statusCode := withStatusCode.StatusCode()
		w.WriteHeader(statusCode)

		if statusCode == http.StatusNoContent {
			return
		}
	}

	_ = json.NewEncoder(w).Encode(rsp)
}

type withStatusCreated struct{}

func (withStatusCreated) StatusCode() int {
	return http.StatusCreated
}
