package buble

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// User is our domain model.
type User struct {
	ID int `json:"id"`
}

func Test_Handler(t *testing.T) {
	resp := httptest.NewRecorder()
	h := &Handler{
		HandlerFunc: HandlerFunc(func(resp *Response, req *Request) {
			resp.SetStatus(200)
			resp.Present(&User{ID: 1})
		}),
	}

	req, _ := http.NewRequest("GET", "", nil)
	h.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Error("Expected 200 OK")
	}

	if resp.Body.String() != "{\"id\":1}\n" {
		t.Error("Expected JSON body")
	}
}
