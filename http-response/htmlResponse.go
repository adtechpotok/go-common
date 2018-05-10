package http_response

import (
	"fmt"
	"net/http"
)

// Struct to html response
type HTMLResponse struct {
	BaseResponse
}

// RenderHtml convert the response to HTML and send it to client
func (res *HTMLResponse) Render(w http.ResponseWriter) {
	w.WriteHeader(res.status)
	w.Write([]byte(fmt.Sprintf("%s %d", res.Body, res.status)))
}
