package http_response

import (
	"net/http"
)

//Server response structure
type BaseResponse struct {
	Body   string `json:"message"`
	status int
}

//CustomErrorState turn on CustomErrorState of the response
func (res *BaseResponse) CustomErrorState(message string, code int) {
	res.setMessage(message)
	res.setStatus(code)
}

//ServerErrorState turn on ServerErrorState of the response
func (res *BaseResponse) ServerErrorState() {
	res.setMessage(http.StatusText(http.StatusInternalServerError))
	res.setStatus(http.StatusInternalServerError)
}

//BadRequestErrorState turn on BadRequestErrorState of the response
func (res *BaseResponse) BadRequestErrorState() {
	res.setMessage(http.StatusText(http.StatusBadRequest))
	res.setStatus(http.StatusBadRequest)
}

//ForbiddenErrorState turn on ForbiddenErrorState of the response
func (res *BaseResponse) ForbiddenErrorState() {
	res.setMessage(http.StatusText(http.StatusForbidden))
	res.setStatus(http.StatusForbidden)
}

//TooManyRequestsErrorState turn on TooManyRequestsErrorState of the response
func (res *BaseResponse) TooManyRequestsErrorState() {
	res.setMessage(http.StatusText(http.StatusTooManyRequests))
	res.setStatus(http.StatusTooManyRequests)
}

//NotFoundErrorState turn on NotFoundErrorState of the response
func (res *BaseResponse) NotFoundErrorState() {
	res.setMessage(http.StatusText(http.StatusNotFound))
	res.setStatus(http.StatusNotFound)
}

//setMessage set response message
func (res *BaseResponse) setMessage(message string) {
	res.Body = message
}

//setStatus set response status
func (res *BaseResponse) setStatus(status int) {
	res.status = status
}
