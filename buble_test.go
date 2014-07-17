package buble

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// model is some domain model.
type model struct {
	ID int
}

// form implements the Form interface.
type form struct {
	body string
}

func (f *form) Validate() error {
	return nil
}

// resource implements the Resource interface.
type resource struct {
	*model
}

func (r *resource) Present() interface{} {
	return struct {
		ID int `json:"int"`
	}{
		ID: r.ID,
	}
}

func Test_Handler(t *testing.T) {
	resp := httptest.NewRecorder()
	h := &Handler{
		Form:     &form{},
		Resource: &resource{},
		HandlerFunc: HandlerFunc(func(resp *Response, req *Request) {
		}),
	}

	req, _ := http.NewRequest("GET", "", nil)
	h.ServeHTTP(resp, req)
}
