package http_response

import (
	`net/http`
	`github.com/json-iterator/go`
)

// JSON response for raw data
type JSONRawResponse struct {
	body   interface{}
	status int
}

// Render render raw data to response
func (response *JSONRawResponse) Render(w http.ResponseWriter) {
	js, err := jsoniter.Marshal(response.body)
	if err != nil {
		panic(Error{err, response})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(response.status)
	w.Write(js)
}
