package http_response

import (
	`net/http`
	`github.com/json-iterator/go`
)

// JSONRaw return empty JSONRawResponse struct
func JSONRaw() *JSONRawResponse {
	return &JSONRawResponse{}
}

// JSON response for raw data
type JSONRawResponse struct {
	body   interface{}
	status int
}

// Render render raw data to response
func (res *JSONRawResponse) Render(w http.ResponseWriter) {
	js, err := jsoniter.Marshal(res.body)
	if err != nil {
		panic(Error{err, res})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(res.status)
	w.Write(js)
}

// setStatus set response Status
func (res *JSONRawResponse) setStatus(status int) Response {
	res.status = status
	return res
}

// setMessage set response message
func (res *JSONRawResponse) setMessage(body interface{}) Response {
	res.body = body
	return res
}

// CustomState turn on CustomState of the response
func (res *JSONRawResponse) CustomState(message interface{}, code int) Response {
	return CustomState(res, message, code)
}

// ServerErrorState turn on ServerErrorState of the response
func (res *JSONRawResponse) ServerErrorState() Response {
	return ServerErrorState(res)
}

// BadRequestErrorState turn on BadRequestErrorState of the response
func (res *JSONRawResponse) BadRequestErrorState() Response {
	return BadRequestErrorState(res)
}

// ForbiddenErrorState turn on ForbiddenErrorState of the response
func (res *JSONRawResponse) ForbiddenErrorState() Response {
	return ForbiddenErrorState(res)
}

// TooManyRequestsErrorState turn on TooManyRequestsErrorState of the response
func (res *JSONRawResponse) TooManyRequestsErrorState() Response {
	return TooManyRequestsErrorState(res)
}

// NotFoundErrorState turn on NotFoundErrorState of the response
func (res *JSONRawResponse) NotFoundErrorState() Response {
	return NotFoundErrorState(res)
}
