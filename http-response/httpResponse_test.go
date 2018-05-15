package http_response

import (
	`testing`
	`net/http`
	`net/http/httptest`
)

func TestErrorStates(t *testing.T) {
	html := HTML()
	json := JSON()
	jsonRaw := JSONRaw()

	html.ServerErrorState()
	json.ServerErrorState()
	jsonRaw.ServerErrorState()

	if html.Body != http.StatusText(http.StatusInternalServerError) && html.status != http.StatusInternalServerError {
		t.Error("ServerErorState html error")
	}

	if json.Body != http.StatusText(http.StatusInternalServerError) && json.status != http.StatusInternalServerError {
		t.Error("ServerErorState json error")
	}

	if jsonRaw.body.(string) != http.StatusText(http.StatusInternalServerError) && jsonRaw.status != http.StatusInternalServerError {
		t.Error("ServerErorState jsonRaw error")
	}

	html.BadRequestErrorState()
	json.BadRequestErrorState()
	jsonRaw.BadRequestErrorState()

	if html.Body != http.StatusText(http.StatusBadRequest) && html.status != http.StatusBadRequest {
		t.Error("BadRequestErrorState html error")
	}

	if json.Body != http.StatusText(http.StatusBadRequest) && json.status != http.StatusBadRequest {
		t.Error("BadRequestErrorState json error")
	}

	if jsonRaw.body.(string) != http.StatusText(http.StatusBadRequest) && jsonRaw.status != http.StatusBadRequest {
		t.Error("BadRequestErrorState jsonRaw error")
	}

	html.ForbiddenErrorState()
	json.ForbiddenErrorState()
	jsonRaw.ForbiddenErrorState()

	if html.Body != http.StatusText(http.StatusForbidden) && html.status != http.StatusForbidden {
		t.Error("ForbiddenErrorState html error")
	}

	if json.Body != http.StatusText(http.StatusForbidden) && json.status != http.StatusForbidden {
		t.Error("ForbiddenErrorState json error")
	}

	if jsonRaw.body.(string) != http.StatusText(http.StatusForbidden) && jsonRaw.status != http.StatusForbidden {
		t.Error("ForbiddenErrorState jsonRaw error")
	}

	html.TooManyRequestsErrorState()
	json.TooManyRequestsErrorState()
	jsonRaw.TooManyRequestsErrorState()

	if html.Body != http.StatusText(http.StatusTooManyRequests) && html.status != http.StatusTooManyRequests {
		t.Error("TooManyRequestsErrorState html error")
	}

	if json.Body != http.StatusText(http.StatusTooManyRequests) && json.status != http.StatusTooManyRequests {
		t.Error("TooManyRequestsErrorState json error")
	}

	if jsonRaw.body.(string) != http.StatusText(http.StatusTooManyRequests) && jsonRaw.status != http.StatusTooManyRequests {
		t.Error("TooManyRequestsErrorState jsonRaw error")
	}

	html.NotFoundErrorState()
	json.NotFoundErrorState()
	jsonRaw.NotFoundErrorState()

	if html.Body != http.StatusText(http.StatusNotFound) && html.status != http.StatusNotFound {
		t.Error("NotFoundErrorState html error")
	}

	if json.Body != http.StatusText(http.StatusNotFound) && json.status != http.StatusNotFound {
		t.Error("NotFoundErrorState json error")
	}

	if jsonRaw.body.(string) != http.StatusText(http.StatusNotFound) && jsonRaw.status != http.StatusNotFound {
		t.Error("NotFoundErrorState jsonRaw error")
	}
}

func TestCustomErrorState(t *testing.T) {
	html := HTML()
	json := JSON()
	jsonRaw := JSONRaw()

	message := "Test"
	status := 100
	html.CustomState(message, status)
	json.CustomState(message, status)
	jsonRaw.CustomState(message, status)

	if html.Body != message && html.status != status {
		t.Error("ServerErorState html error")
	}

	if json.Body != message && json.status != status {
		t.Error("ServerErorState json error")
	}

	if jsonRaw.body.(string) != message && jsonRaw.status != status {
		t.Error("ServerErorState jsonRaw error")
	}
}

func TestHTMLResponse_Render(t *testing.T) {
	w := httptest.ResponseRecorder{}
	HTML().NotFoundErrorState().Render(&w)
	if w.Code != http.StatusNotFound {
		t.Error("Html render problem")
	}
}

func TestJSONResponse_Render(t *testing.T) {
	w := httptest.ResponseRecorder{}
	JSON().NotFoundErrorState().Render(&w)
	if w.Code != http.StatusNotFound {
		t.Error("Json render problem")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Trash value to json work. JSON")
		}
	}()
	value := make(chan int)
	JSON().CustomState(value, 100).Render(&w)

}

func TestJSONRawResponse_Render(t *testing.T) {
	w := httptest.ResponseRecorder{}
	JSONRaw().NotFoundErrorState().Render(&w)
	if w.Code != http.StatusNotFound {
		t.Error("JSONRaw render problem")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Trash value to json work. JsonRaw")
		}
	}()
	value := make(chan int)
	JSONRaw().CustomState(value, 100).Render(&w)
}
