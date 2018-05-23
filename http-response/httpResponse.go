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

// setCustomResult turn on setCustomResult state of the response
func setCustomResult(res Response, message interface{}, code int) Response {
	res.setMessage(message)
	res.setStatus(code)

	return res
}

// setSuccessResult turn on setSuccessResult state of the response
func setSuccessResult(res Response, message interface{}) Response {
	res.setMessage(message)
	res.setStatus(http.StatusOK)

	return res
}

// setServerError turn on setServerError state of the response
func setServerError(res Response) Response {
	res.setMessage(http.StatusText(http.StatusInternalServerError))
	res.setStatus(http.StatusInternalServerError)

	return res
}

// setBadRequestError turn on setBadRequestError state of the response
func setBadRequestError(res Response) Response {
	res.setMessage(http.StatusText(http.StatusBadRequest))
	res.setStatus(http.StatusBadRequest)

	return res
}

// setForbiddenError turn on setForbiddenError state of the response
func setForbiddenError(res Response) Response {
	res.setMessage(http.StatusText(http.StatusForbidden))
	res.setStatus(http.StatusForbidden)

	return res
}

// setTooManyRequestsError turn on setTooManyRequestsError state of the response
func setTooManyRequestsError(res Response) Response {
	res.setMessage(http.StatusText(http.StatusTooManyRequests))
	res.setStatus(http.StatusTooManyRequests)

	return res
}

// setNotFoundError turn on setNotFoundError state of the response
func setNotFoundError(res Response) Response {
	res.setMessage(http.StatusText(http.StatusNotFound))
	res.setStatus(http.StatusNotFound)

	return res
}
