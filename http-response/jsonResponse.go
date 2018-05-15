package http_response

import (
	`net/http`
	`github.com/json-iterator/go`
)

// JSON return empty JSONResponse struct
func JSON() *JSONResponse {
	return &JSONResponse{}
}

// JSON response
type JSONResponse struct {
	Body   string `json:"message"`
	status int
}

// RenderJson convert the response to JSON and send it to client
func (res *JSONResponse) Render(w http.ResponseWriter) {
	js, err := jsoniter.Marshal(res)
	if err != nil {
		panic(Error{err, res})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(res.status)
	w.Write(js)
}

// setStatus set response Status
func (res *JSONResponse) setStatus(status int) Response {
	res.status = status
	return res
}

// setMessage set response message
func (res *JSONResponse) setMessage(message interface{}) Response {
	res.Body = message.(string)
	return res
}

// CustomState turn on CustomState of the response
func (res *JSONResponse) CustomState(message interface{}, code int) Response {
	return CustomState(res, message, code)
}

// ServerErrorState turn on ServerErrorState of the response
func (res *JSONResponse) ServerErrorState() Response {
	return ServerErrorState(res)
}

// BadRequestErrorState turn on BadRequestErrorState of the response
func (res *JSONResponse) BadRequestErrorState() Response {
	return BadRequestErrorState(res)
}

// ForbiddenErrorState turn on ForbiddenErrorState of the response
func (res *JSONResponse) ForbiddenErrorState() Response {
	return ForbiddenErrorState(res)
}

// TooManyRequestsErrorState turn on TooManyRequestsErrorState of the response
func (res *JSONResponse) TooManyRequestsErrorState() Response {
	return TooManyRequestsErrorState(res)
}

// NotFoundErrorState turn on NotFoundErrorState of the response
func (res *JSONResponse) NotFoundErrorState() Response {
	return NotFoundErrorState(res)
}
