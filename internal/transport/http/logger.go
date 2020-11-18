package http

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type loggerInterface interface {
	Println(v ...interface{})
}

type loggerMW struct {
	next   http.Handler
	logger loggerInterface
}

func logger(logger loggerInterface) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return &loggerMW{
			next:   next,
			logger: logger,
		}
	}
}

func dumpBody(ctp string) bool {
	allowedTypes := []string{
		"application/json",
		"application/problem+json",
	}

	for i := range allowedTypes {
		if allowedTypes[i] == ctp {
			return true
		}
	}

	return false
}

func (mw *loggerMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// ignore body if is file
	reqDump, err := httputil.DumpRequest(r, dumpBody(r.Header.Get("Content-Type")))
	if err != nil {
		mw.logger.Println(errors.Wrap(errors.WithStack(err), "error dumping request"))
	} else {
		reqParts := strings.Split(string(reqDump), "\n")
		req := ""
		for i := range reqParts {
			req += fmt.Sprintf("%s\n", aurora.BrightCyan(reqParts[i]))
		}
		mw.logger.Println(req)
	}

	crw := &customRW{rw: w, statusCode: http.StatusOK, body: bytes.NewBufferString("")}
	w = crw

	mw.next.ServeHTTP(w, r)

	res := aurora.Cyan(fmt.Sprintf("%s %s\n\n", r.Method, r.RequestURI)).String()
	res += aurora.Cyan(fmt.Sprintf("%s %d %s\n", r.Proto, crw.statusCode, http.StatusText(crw.statusCode))).String()
	for k, vv := range w.Header() {
		for i := range vv {
			res += aurora.Cyan(fmt.Sprintf("%s: %s\n", k, vv[i])).String()
		}
	}
	if dumpBody(w.Header().Get("Content-Type")) {
		res += fmt.Sprintf("\n%s\n\n", aurora.Cyan(crw.body.String()))
	}
	res += aurora.Cyan(fmt.Sprintf(
		"Response code: %d (%s); Time: %s; Content length: %d bytes",
		crw.statusCode,
		http.StatusText(crw.statusCode),
		time.Since(start),
		crw.body.Len(),
	)).String()
	mw.logger.Println(res)
}

type customRW struct {
	rw         http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (crw *customRW) Header() http.Header {
	return crw.rw.Header()
}

func (crw *customRW) Write(i []byte) (int, error) {
	crw.body.Write(i)
	return crw.rw.Write(i)
}

func (crw *customRW) WriteHeader(statusCode int) {
	crw.statusCode = statusCode
	crw.rw.WriteHeader(statusCode)
}
