package http_response

import (
	"net/http"
	"github.com/json-iterator/go"
)

// Struct to JSON response
type JSONResponse struct {
	BaseResponse
}

//RenderJson convert the response to JSON and send it to client
func (res JSONResponse) Render(w http.ResponseWriter) {
	js, err := jsoniter.Marshal(res)
	if err != nil {
		panic(Error{err, &res})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(res.status)
	w.Write(js)
}
