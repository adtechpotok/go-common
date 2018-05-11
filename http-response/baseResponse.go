package http_response

import (
	"net/http"
)

const jsonType = "JSON"
const jsonRawType = "JSONRaw"


// Server response structure
type Response struct {
	renderType   string
	Body   string `json:"message"`
	status int
}

// CustomErrorState turn on CustomErrorState of the response
func (res *Response) CustomErrorState(message string, code int) Response {
	res.setMessage(message)
	res.setStatus(code)

	return res
}

// ServerErrorState turn on ServerErrorState of the response
func (res *Response) ServerErrorState() *Response {
	res.setMessage(http.StatusText(http.StatusInternalServerError))
	res.setStatus(http.StatusInternalServerError)

	return res
}

// BadRequestErrorState turn on BadRequestErrorState of the response
func (res *Response) BadRequestErrorState() *Response {
	res.setMessage(http.StatusText(http.StatusBadRequest))
	res.setStatus(http.StatusBadRequest)

	return res
}

// ForbiddenErrorState turn on ForbiddenErrorState of the response
func (res *Response) ForbiddenErrorState() *Response {
	res.setMessage(http.StatusText(http.StatusForbidden))
	res.setStatus(http.StatusForbidden)

	return res
}

// TooManyRequestsErrorState turn on TooManyRequestsErrorState of the response
func (res *Response) TooManyRequestsErrorState() *Response {
	res.setMessage(http.StatusText(http.StatusTooManyRequests))
	res.setStatus(http.StatusTooManyRequests)

	return res
}

// NotFoundErrorState turn on NotFoundErrorState of the response
func (res *Response) NotFoundErrorState() *Response {
	res.setMessage(http.StatusText(http.StatusNotFound))
	res.setStatus(http.StatusNotFound)

	return res
}

// setMessage set response message
func (res *Response) setMessage(message string) *Response {
	res.Body = message

	return res
}

// setStatus set response Status
func (res *Response) setStatus(status int) *Response {
	res.status = status

	return res
}
