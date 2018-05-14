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
func (res *HTMLResponse) Render(w http.ResponseWriter) {
	w.WriteHeader(res.status)
	w.Write([]byte(fmt.Sprintf("%s %d", res.Body, res.status)))
}

// setStatus set response Status
func (res *HTMLResponse) setStatus(status int) Response {
	res.status = status
	return res
}

// setMessage set response message
func (res *HTMLResponse) setMessage(message interface{}) Response {
	res.Body = message.(string)
	return res
}

// CustomState turn on CustomState of the response
func (res *HTMLResponse) CustomState(message interface{}, code int) Response {
	return CustomState(res, message, code)
}

// ServerErrorState turn on ServerErrorState of the response
func (res *HTMLResponse) ServerErrorState() Response {
	return ServerErrorState(res)
}

// BadRequestErrorState turn on BadRequestErrorState of the response
func (res *HTMLResponse) BadRequestErrorState() Response {
	return BadRequestErrorState(res)
}

// ForbiddenErrorState turn on ForbiddenErrorState of the response
func (res *HTMLResponse) ForbiddenErrorState() Response {
	return ForbiddenErrorState(res)
}

// TooManyRequestsErrorState turn on TooManyRequestsErrorState of the response
func (res *HTMLResponse) TooManyRequestsErrorState() Response {
	return TooManyRequestsErrorState(res)
}

// NotFoundErrorState turn on NotFoundErrorState of the response
func (res *HTMLResponse) NotFoundErrorState() Response {
	return NotFoundErrorState(res)
}
