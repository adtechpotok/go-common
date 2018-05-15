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
	CustomResult(message interface{}, code int) Response
	SuccessResult(body interface{}) Response
	ServerError() Response
	BadRequestError() Response
	ForbiddenError() Response
	TooManyRequestsError() Response
	NotFoundError() Response
	setMessage(body interface{}) Response
	setStatus(status int) Response
}

// custom turn on custom state of the response
func custom(res Response, message interface{}, code int) Response {
	res.setMessage(message)
	res.setStatus(code)

	return res
}

// success turn on success state of the response
func success(res Response, message interface{}) Response {
	res.setMessage(message)
	res.setStatus(http.StatusOK)

	return res
}

// serverError turn on serverError state of the response
func serverError(res Response) Response {
	res.setMessage(http.StatusText(http.StatusInternalServerError))
	res.setStatus(http.StatusInternalServerError)

	return res
}

// badRequestError turn on badRequestError state of the response
func badRequestError(res Response) Response {
	res.setMessage(http.StatusText(http.StatusBadRequest))
	res.setStatus(http.StatusBadRequest)

	return res
}

// forbiddenError turn on forbiddenError state of the response
func forbiddenError(res Response) Response {
	res.setMessage(http.StatusText(http.StatusForbidden))
	res.setStatus(http.StatusForbidden)

	return res
}

// tooManyRequestsError turn on tooManyRequestsError state of the response
func tooManyRequestsError(res Response) Response {
	res.setMessage(http.StatusText(http.StatusTooManyRequests))
	res.setStatus(http.StatusTooManyRequests)

	return res
}

// notFoundError turn on notFoundError state of the response
func notFoundError(res Response) Response {
	res.setMessage(http.StatusText(http.StatusNotFound))
	res.setStatus(http.StatusNotFound)

	return res
}
