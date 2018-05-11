package http_response

import (
	`net/http`
	`github.com/json-iterator/go`
)

/*
// JSON response for raw data
type JSONRawResponse struct {
	Body   interface{}
	Status int
}

// Render render raw data to response
func (response *JSONRawResponse) Render(w http.ResponseWriter) {
	js, err := jsoniter.Marshal(response.Body)
	if err != nil {
		panic(Error{err, response})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(response.Status)
	w.Write(js)
}
*/
func JSONRawRender(response Response, w http.ResponseWriter) {
	js, err := jsoniter.Marshal(response.Body)
	if err != nil {
		panic(Error{err, response})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(response.status)
	w.Write(js)
}

func JSONRaw() Response {
	return Response{renderType: jsonRawType}
}
