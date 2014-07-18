package buble

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// User is our domain model.
type User struct {
	ID int
}

// UserCreateForm implements the Form interface.
type UserCreateForm struct {
	Name string
}

func (f *UserCreateForm) Validate() error {
	return nil
}

// UserResource implements the Resource interface.
type UserResource struct {
	*User
}

func (r *UserResource) Present() interface{} {
	return struct {
		ID int `json:"int"`
	}{
		ID: r.ID,
	}
}

func Test_Handler(t *testing.T) {
	resp := httptest.NewRecorder()
	h := &Handler{
		Form:     &UserCreateForm{},
		Resource: &UserResource{},
		HandlerFunc: HandlerFunc(func(resp *Response, req *Request) {
		}),
	}

	req, _ := http.NewRequest("GET", "", nil)
	h.ServeHTTP(resp, req)
}
