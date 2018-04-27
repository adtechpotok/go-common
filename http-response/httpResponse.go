package http_response

import "net/http"

type Error struct{
	Error    error
	Response Response
}


type Response interface {
	Render(w http.ResponseWriter)
	CustomErrorState(message string, code int)
	ServerErrorState()
	BadRequestErrorState()
	ForbiddenErrorState()
	TooManyRequestsErrorState()
	NotFoundErrorState()
	setMessage(message string)
	setStatus(status int)
}
