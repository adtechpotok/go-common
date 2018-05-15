package http_response

import (
	`net/http`
	`github.com/json-iterator/go`
)

// JSONMessage return empty JSONMessageResponse struct
func JSONMessage() *JSONMessageResponse {
	return &JSONMessageResponse{}
}

// JSON response with string body
type JSONMessageResponse struct {
	Body   string `json:"message"`
	status int
}

// RenderJson convert the response to JSONMessage and send it to client
func (m *JSONMessageResponse) Render(w http.ResponseWriter) {
	js, _ := jsoniter.Marshal(m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(m.status)
	w.Write(js)
}

// setStatus set response Status
func (m *JSONMessageResponse) setStatus(status int) Response {
	m.status = status
	return m
}

// setMessage set response message
func (m *JSONMessageResponse) setMessage(message interface{}) Response {
	m.Body = message.(string)
	return m
}

// CustomResult turn on custom state of the response
func (m *JSONMessageResponse) CustomResult(message interface{}, code int) Response {
	return custom(m, message, code)
}

// ServerError turn on ServerError state of the response
func (m *JSONMessageResponse) ServerError() Response {
	return serverError(m)
}

// BadRequestError turn on BadRequestError state of the response
func (m *JSONMessageResponse) BadRequestError() Response {
	return badRequestError(m)
}

// ForbiddenError turn on ForbiddenError state of the response
func (m *JSONMessageResponse) ForbiddenError() Response {
	return forbiddenError(m)
}

// TooManyRequestsError turn on TooManyRequestsError state of the response
func (m *JSONMessageResponse) TooManyRequestsError() Response {
	return tooManyRequestsError(m)
}

// NotFoundError turn on NotFoundError state of the response
func (m *JSONMessageResponse) NotFoundError() Response {
	return notFoundError(m)
}

// SuccessResult turn on success state of the response
func (m *JSONMessageResponse) SuccessResult(body interface{}) Response {
	return success(m, body)
}
