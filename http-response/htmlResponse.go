package http_response

import (
	`net/http`
	`fmt`
)

// HTML return empty HTMLResponse struct
func HTML() *HTMLResponse {
	return &HTMLResponse{}
}

// HTML response
type HTMLResponse struct {
	Body   string `json:"message"`
	status int
}

// RenderHtml convert the response to HTML and send it to client
func (m *HTMLResponse) Render(w http.ResponseWriter) {
	w.WriteHeader(m.status)
	w.Write([]byte(fmt.Sprintf("%s %d", m.Body, m.status)))
}

// setStatus set response Status
func (m *HTMLResponse) setStatus(status int) Response {
	m.status = status
	return m
}

// setMessage set response message
func (m *HTMLResponse) setMessage(message interface{}) Response {
	m.Body = message.(string)
	return m
}

// Custom turn on custom state of the response
func (m *HTMLResponse) Custom(message interface{}, code int) Response {
	return custom(m, message, code)
}

// ServerError turn on ServerError state of the response
func (m *HTMLResponse) ServerError() Response {
	return serverError(m)
}

// BadRequestError turn on BadRequestError state of the response
func (m *HTMLResponse) BadRequestError() Response {
	return badRequestError(m)
}

// ForbiddenError turn on ForbiddenError state of the response
func (m *HTMLResponse) ForbiddenError() Response {
	return forbiddenError(m)
}

// TooManyRequestsError turn on TooManyRequestsError state of the response
func (m *HTMLResponse) TooManyRequestsError() Response {
	return tooManyRequestsError(m)
}

// NotFoundError turn on NotFoundError state of the response
func (m *HTMLResponse) NotFoundError() Response {
	return notFoundError(m)
}

// Success turn on success state of the response
func (m *HTMLResponse) Success(body interface{}) Response {
	return success(m, body)
}
