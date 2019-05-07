package api

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func newErrResponse(status int, err error) *ErrResponse {
	statusText := http.StatusText(status)
	if statusText == "" {
		log.Fatalf("unknown http status %d", status)
	}

	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: status,
		StatusText:     statusText,
		ErrorText:      err.Error(),
	}
}

func ErrInvalidRequest(err error) render.Renderer {
	return newErrResponse(http.StatusBadRequest, err)
}

func ErrRender(err error) render.Renderer {
	return newErrResponse(http.StatusUnprocessableEntity, err)
}

func ErrInternalError(err error) render.Renderer {
	return newErrResponse(http.StatusInternalServerError, err)
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: http.StatusText(http.StatusNotFound)}

// ErrNotImplemented is not used currently (as everything is implemented), but it's supposed to be returned whenever
// something was not implemented.
var ErrNotImplemented = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: "This is not implemented, " +
	"as I don't want to play with mongo more than I have to, but you get the idea. " +
	"Please check the in-memory implementation."}
