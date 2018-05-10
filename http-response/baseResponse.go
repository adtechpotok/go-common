package http_response

import (
	"net/http"
)

// Server response structure
type BaseResponse struct {
	Body   string `json:"message"`
	status int
}

// CustomErrorState turn on CustomErrorState of the response
func (res *BaseResponse) CustomErrorState(message string, code int) *BaseResponse {
	res.setMessage(message)
	res.setStatus(code)

	return res
}

// ServerErrorState turn on ServerErrorState of the response
func (res *BaseResponse) ServerErrorState() *BaseResponse {
	res.setMessage(http.StatusText(http.StatusInternalServerError))
	res.setStatus(http.StatusInternalServerError)

	return res
}

// BadRequestErrorState turn on BadRequestErrorState of the response
func (res *BaseResponse) BadRequestErrorState() *BaseResponse {
	res.setMessage(http.StatusText(http.StatusBadRequest))
	res.setStatus(http.StatusBadRequest)

	return res
}

// ForbiddenErrorState turn on ForbiddenErrorState of the response
func (res *BaseResponse) ForbiddenErrorState() *BaseResponse {
	res.setMessage(http.StatusText(http.StatusForbidden))
	res.setStatus(http.StatusForbidden)

	return res
}

// TooManyRequestsErrorState turn on TooManyRequestsErrorState of the response
func (res *BaseResponse) TooManyRequestsErrorState() *BaseResponse {
	res.setMessage(http.StatusText(http.StatusTooManyRequests))
	res.setStatus(http.StatusTooManyRequests)

	return res
}

// NotFoundErrorState turn on NotFoundErrorState of the response
func (res *BaseResponse) NotFoundErrorState() *BaseResponse {
	res.setMessage(http.StatusText(http.StatusNotFound))
	res.setStatus(http.StatusNotFound)

	return res
}

// setMessage set response message
func (res *BaseResponse) setMessage(message string) *BaseResponse {
	res.Body = message

	return res
}

// setStatus set response status
func (res *BaseResponse) setStatus(status int) *BaseResponse {
	res.status = status

	return res
}
