package http_response

import (
	`net/http`
)

type Error struct {
	Error    error
	Response Response
}

// Regular response
type Response interface {
	Render(w http.ResponseWriter)
	CustomState(message interface{}, code int) Response
	ServerErrorState() Response
	BadRequestErrorState() Response
	ForbiddenErrorState() Response
	TooManyRequestsErrorState() Response
	NotFoundErrorState() Response
	setMessage(body interface{}) Response
	setStatus(status int) Response
}

// CustomState turn on CustomState of the response
func CustomState(res Response, message interface{}, code int) Response {
	res.setMessage(message)
	res.setStatus(code)

	return res
}

// ServerErrorState turn on ServerErrorState of the response
func ServerErrorState(res Response) Response {
	res.setMessage(http.StatusText(http.StatusInternalServerError))
	res.setStatus(http.StatusInternalServerError)

	return res
}

// BadRequestErrorState turn on BadRequestErrorState of the response
func BadRequestErrorState(res Response) Response {
	res.setMessage(http.StatusText(http.StatusBadRequest))
	res.setStatus(http.StatusBadRequest)

	return res
}

// ForbiddenErrorState turn on ForbiddenErrorState of the response
func ForbiddenErrorState(res Response) Response {
	res.setMessage(http.StatusText(http.StatusForbidden))
	res.setStatus(http.StatusForbidden)

	return res
}

// TooManyRequestsErrorState turn on TooManyRequestsErrorState of the response
func TooManyRequestsErrorState(res Response) Response {
	res.setMessage(http.StatusText(http.StatusTooManyRequests))
	res.setStatus(http.StatusTooManyRequests)

	return res
}

// NotFoundErrorState turn on NotFoundErrorState of the response
func NotFoundErrorState(res Response) Response {
	res.setMessage(http.StatusText(http.StatusNotFound))
	res.setStatus(http.StatusNotFound)

	return res
}
