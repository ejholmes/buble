package buble

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// User is our domain model.
type User struct {
	ID int `json:"id"`
}

// Form object
type UserCreateForm struct {
	Name string `json:"name"`
}

func Test_Handler(t *testing.T) {
	type expectation struct {
		status      int
		body        string
		contentType string
	}

	tests := []struct {
		method   string
		body     string
		fn       HandlerFunc
		expected expectation
	}{
		{
			fn: HandlerFunc(func(w ResponseWriter, r *Request) {
				w.SetStatus(200)
				w.Present(&User{ID: 1})
			}),
			expected: expectation{
				status:      200,
				body:        `{"id":1}` + "\n",
				contentType: "application/json",
			},
		},
		{
			fn: HandlerFunc(func(w ResponseWriter, r *Request) {
				w.SetStatus(200)
			}),
			expected: expectation{
				status:      200,
				body:        `{}` + "\n",
				contentType: "application/json",
			},
		},
		{
			method: "POST",
			body:   `{"name":"Eric Holmes"}`,
			fn: HandlerFunc(func(w ResponseWriter, r *Request) {
				var f UserCreateForm
				r.Decode(&f)

				w.SetStatus(200)
				w.Present(f)
			}),
			expected: expectation{
				status:      200,
				body:        `{"name":"Eric Holmes"}` + "\n",
				contentType: "application/json",
			},
		},
	}

	for _, test := range tests {
		method := test.method
		if method == "" {
			method = "GET"
		}

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(method, "", bytes.NewReader([]byte(test.body)))

		h := &Handler{HandlerFunc: test.fn}
		h.ServeHTTP(resp, req)

		if resp.Code != test.expected.status {
			t.Errorf("Status: Want %v, Got %v", test.expected.status, resp.Code)
		}

		if resp.Body.String() != test.expected.body {
			t.Errorf("Body: Want %v, Got %v", test.expected.body, resp.Body.String())
		}

		contentType := resp.HeaderMap.Get("Content-Type")
		if contentType != test.expected.contentType {
			t.Errorf("Content-Type: Want %v, Got %v", test.expected.contentType, contentType)
		}
	}
}
