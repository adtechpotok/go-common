package http_response

import (
	"github.com/json-iterator/go"
	"net/http"
)

/*
// Struct to JSON response
type JSONResponse struct {
	Response
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
*/
func JSONRender(response Response, w http.ResponseWriter) {
	js, err := jsoniter.Marshal(response)
	if err != nil {
		panic(Error{err, response})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(response.status)
	w.Write(js)
}

func JSON() Response {
	return Response{renderType: jsonType}
}
