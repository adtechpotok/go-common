package http_response

import (
	`testing`
	`net/http`
	`net/http/httptest`
)

func TestErrorStates(t *testing.T) {
	html := HTML()
	jsonMessage := JSONMessage()
	json := JSON()

	html.ServerError()
	jsonMessage.ServerError()
	json.ServerError()

	if html.Body != http.StatusText(http.StatusInternalServerError) && html.status != http.StatusInternalServerError {
		t.Error("ServerErorState html error")
	}

	if jsonMessage.Body != http.StatusText(http.StatusInternalServerError) && jsonMessage.status != http.StatusInternalServerError {
		t.Error("ServerErorState jsonMessage error")
	}

	if json.body.(string) != http.StatusText(http.StatusInternalServerError) && json.status != http.StatusInternalServerError {
		t.Error("ServerErorState jsonMessage error")
	}

	html.BadRequestError()
	jsonMessage.BadRequestError()
	json.BadRequestError()

	if html.Body != http.StatusText(http.StatusBadRequest) && html.status != http.StatusBadRequest {
		t.Error("BadRequestError html error")
	}

	if jsonMessage.Body != http.StatusText(http.StatusBadRequest) && jsonMessage.status != http.StatusBadRequest {
		t.Error("BadRequestError jsonMessage error")
	}

	if json.body.(string) != http.StatusText(http.StatusBadRequest) && json.status != http.StatusBadRequest {
		t.Error("BadRequestError jsonMessage error")
	}

	html.ForbiddenError()
	jsonMessage.ForbiddenError()
	json.ForbiddenError()

	if html.Body != http.StatusText(http.StatusForbidden) && html.status != http.StatusForbidden {
		t.Error("ForbiddenError html error")
	}

	if jsonMessage.Body != http.StatusText(http.StatusForbidden) && jsonMessage.status != http.StatusForbidden {
		t.Error("ForbiddenError jsonMessage error")
	}

	if json.body.(string) != http.StatusText(http.StatusForbidden) && json.status != http.StatusForbidden {
		t.Error("ForbiddenError jsonMessage error")
	}

	html.TooManyRequestsError()
	jsonMessage.TooManyRequestsError()
	json.TooManyRequestsError()

	if html.Body != http.StatusText(http.StatusTooManyRequests) && html.status != http.StatusTooManyRequests {
		t.Error("TooManyRequestsError html error")
	}

	if jsonMessage.Body != http.StatusText(http.StatusTooManyRequests) && jsonMessage.status != http.StatusTooManyRequests {
		t.Error("TooManyRequestsError jsonMessage error")
	}

	if json.body.(string) != http.StatusText(http.StatusTooManyRequests) && json.status != http.StatusTooManyRequests {
		t.Error("TooManyRequestsError jsonMessage error")
	}

	html.NotFoundError()
	jsonMessage.NotFoundError()
	json.NotFoundError()

	if html.Body != http.StatusText(http.StatusNotFound) && html.status != http.StatusNotFound {
		t.Error("NotFoundError html error")
	}

	if jsonMessage.Body != http.StatusText(http.StatusNotFound) && jsonMessage.status != http.StatusNotFound {
		t.Error("NotFoundError jsonMessage error")
	}

	if json.body.(string) != http.StatusText(http.StatusNotFound) && json.status != http.StatusNotFound {
		t.Error("NotFoundError jsonMessage error")
	}
}

func TestCustomErrorState(t *testing.T) {
	html := HTML()
	jsonMessage := JSONMessage()
	json := JSON()

	message := "Test"
	status := 100
	html.Custom(message, status)
	jsonMessage.Custom(message, status)
	json.Custom(message, status)

	if html.Body != message && html.status != status {
		t.Error("ServerErorState html error")
	}

	if jsonMessage.Body != message && jsonMessage.status != status {
		t.Error("ServerErorState jsonMessage error")
	}

	if json.body.(string) != message && json.status != status {
		t.Error("ServerErorState jsonMessage error")
	}
}

func TestHTMLResponse_Render(t *testing.T) {
	w := httptest.ResponseRecorder{}
	HTML().NotFoundError().Render(&w)
	if w.Code != http.StatusNotFound {
		t.Error("Html render problem")
	}
}

func TestJSONResponse_Render(t *testing.T) {
	w := httptest.ResponseRecorder{}
	JSONMessage().NotFoundError().Render(&w)
	if w.Code != http.StatusNotFound {
		t.Error("Json render problem")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Trash value to json work. JSONMessage")
		}
	}()
	value := make(chan int)
	JSONMessage().Custom(value, 100).Render(&w)

}

func TestJSONRawResponse_Render(t *testing.T) {
	w := httptest.ResponseRecorder{}
	JSON().NotFoundError().Render(&w)
	if w.Code != http.StatusNotFound {
		t.Error("JSON render problem")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Trash value to json work. json")
		}
	}()
	value := make(chan int)
	JSON().Custom(value, 100).Render(&w)
}

func TestSuccessResponse(t *testing.T) {
	html := HTML()
	jsonMessage := JSONMessage()
	json := JSON()

	message := "1"

	html.Success(message)
	jsonMessage.Success(message)
	json.Success(message)

	if html.Body != message && html.status != http.StatusOK {
		t.Error("Success html error")
	}

	if jsonMessage.Body != message && jsonMessage.status != http.StatusOK {
		t.Error("Success jsonMessage error")
	}

	if json.body.(string) != message && json.status != http.StatusOK {
		t.Error("Success jsonMessage error")
	}
}
