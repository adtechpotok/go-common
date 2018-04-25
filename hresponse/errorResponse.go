package hresponse

import (
	"net/http"
	"fmt"
)

var ServerError = ErrorResponse{http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError}
var BadRequestError = ErrorResponse{http.StatusText(http.StatusBadRequest), http.StatusBadRequest}
var ForbiddenError = ErrorResponse{http.StatusText(http.StatusForbidden), http.StatusForbidden}
var TooManyRequestsError = ErrorResponse{http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests}
var NotFoundError = ErrorResponse{http.StatusText(http.StatusNotFound), http.StatusNotFound}

type ErrorResponse struct {
	Message string `json:"message"`
	status  int
}

func (res ErrorResponse) RenderHtml(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(res.status)
	w.Write([]byte(fmt.Sprintf("%s %d", res.Message, res.status)))
}

func (res ErrorResponse) RenderJson(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(res.status)
	ResultJson(res, w)
}

func (res ErrorResponse) ServerHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(res.status)
	ResultJson(res, w)
}