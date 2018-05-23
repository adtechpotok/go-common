package http_response

import (
	`net/http`
	`github.com/json-iterator/go`
)

// JSON return empty JSONResponse struct
func JSON() *JSONResponse {
	return &JSONResponse{}
}

// JSON response with interface body
type JSONResponse struct {
	body   interface{}
	status int
}

// Render render struct to response
func (m *JSONResponse) Render(w http.ResponseWriter) {
	js, err := jsoniter.Marshal(m.body)
	if err != nil {
		panic(Error{err, m})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(m.status)
	w.Write(js)
}

// setStatus set response Status
func (m *JSONResponse) setStatus(status int) Response {
	m.status = status
	return m
}

// setMessage set response message
func (m *JSONResponse) setMessage(body interface{}) Response {
	m.body = body
	return m
}

// CustomResult turn on custom state of the response
func (m *JSONResponse) CustomResult(message interface{}, code int) Response {
	return setCustomResult(m, message, code)
}

// ServerError turn on ServerError state of the response
func (m *JSONResponse) ServerError() Response {
	return setServerError(m)
}

// BadRequestError turn on BadRequestError state of the response
func (m *JSONResponse) BadRequestError() Response {
	return setBadRequestError(m)
}

// ForbiddenError turn on ForbiddenError state of the response
func (m *JSONResponse) ForbiddenError() Response {
	return setForbiddenError(m)
}

// TooManyRequestsError turn on TooManyRequestsError state of the response
func (m *JSONResponse) TooManyRequestsError() Response {
	return setTooManyRequestsError(m)
}

// NotFoundError turn on NotFoundError state of the response
func (m *JSONResponse) NotFoundError() Response {
	return setNotFoundError(m)
}

// SuccessResult turn on success state of the response
func (m *JSONResponse) SuccessResult(body interface{}) Response {
	return setSuccessResult(m, body)
}
