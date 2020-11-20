package http

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"net/http"
)

type Problem struct {
	Type       string
	Title      string
	Status     int
	Detail     string
	Instance   string
	Extensions map[string]interface{}
}

func (p Problem) MarshalJSON() ([]byte, error) {
	c := make(map[string]interface{})
	c["type"] = "about:blank"
	if p.Type != "" {
		c["type"] = p.Type
	}

	c["status"] = http.StatusInternalServerError
	if p.Status != 0 {
		c["status"] = p.Status
	}

	c["title"] = http.StatusText(c["status"].(int))
	if p.Title != "" {
		c["title"] = p.Title
	}

	c["detail"] = p.Detail

	if p.Instance != "" {
		c["instance"] = p.Instance
	}

	for k, v := range p.Extensions {
		switch k {
		case "type", "status", "title", "detail", "instance":
			c["_"+k] = v
		default:
			c[k] = v
		}
	}

	return json.Marshal(c)
}

func (p Problem) StatusCode() int {
	if p.Status == 0 {
		return http.StatusInternalServerError
	}

	return p.Status
}

func (p Problem) Header() http.Header {
	res := make(http.Header)
	res.Set("Content-Type", "application/problem+json")

	return res
}

type ProblemOption func(e *Problem)

func setExtension(key string, val interface{}) ProblemOption {
	return func(e *Problem) {
		if e.Extensions == nil {
			e.Extensions = make(map[string]interface{})
		}
		e.Extensions[key] = val
	}
}

func internalServerError(err error, options ...ProblemOption) Problem {
	log.Println(fmt.Sprintf("%+v", err))

	id := sentry.CaptureException(err)

	e := Problem{
		Status: http.StatusInternalServerError,
		Extensions: map[string]interface{}{
			"tracking_code": id,
		},
	}

	for i := range options {
		options[i](&e)
	}

	return e
}

func badRequest(detail string, options ...ProblemOption) Problem {
	e := Problem{
		Status:     http.StatusBadRequest,
		Detail:     detail,
		Extensions: map[string]interface{}{},
	}

	for i := range options {
		options[i](&e)
	}

	return e
}

func unauthorized(detail string, options ...ProblemOption) Problem {
	e := Problem{
		Status:     http.StatusUnauthorized,
		Detail:     detail,
		Extensions: map[string]interface{}{},
	}

	for i := range options {
		options[i](&e)
	}

	return e
}

func forbidden(detail string, options ...ProblemOption) Problem {
	e := Problem{
		Status:     http.StatusForbidden,
		Detail:     detail,
		Extensions: map[string]interface{}{},
	}

	for i := range options {
		options[i](&e)
	}

	return e
}

func notFound(detail string, options ...ProblemOption) Problem {
	e := Problem{
		Status:     http.StatusNotFound,
		Detail:     detail,
		Extensions: map[string]interface{}{},
	}

	for i := range options {
		options[i](&e)
	}

	return e
}

func conflict(detail string, options ...ProblemOption) Problem {
	e := Problem{
		Status:     http.StatusConflict,
		Detail:     detail,
		Extensions: map[string]interface{}{},
	}

	for i := range options {
		options[i](&e)
	}

	return e
}

func unsupportedMediaType(detail string, options ...ProblemOption) Problem {
	e := Problem{
		Status:     http.StatusUnsupportedMediaType,
		Detail:     detail,
		Extensions: map[string]interface{}{},
	}

	for i := range options {
		options[i](&e)
	}

	return e
}

func tooManyRequests(detail string, options ...ProblemOption) Problem {
	e := Problem{
		Status:     http.StatusTooManyRequests,
		Detail:     detail,
		Extensions: map[string]interface{}{},
	}

	for i := range options {
		options[i](&e)
	}

	return e
}

func serviceUnavailable(detail string, options ...ProblemOption) Problem {
	e := Problem{
		Status:     http.StatusServiceUnavailable,
		Detail:     detail,
		Extensions: map[string]interface{}{},
	}

	for i := range options {
		options[i](&e)
	}

	return e
}
