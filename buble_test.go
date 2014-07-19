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
		body     string
		fn       HandlerFunc
		expected expectation
	}{
		{
			fn: HandlerFunc(func(resp *Response, req *Request) {
				resp.SetStatus(200)
				resp.Present(&User{ID: 1})
			}),
			expected: expectation{
				status:      200,
				body:        `{"id":1}` + "\n",
				contentType: "application/json",
			},
		},
		{
			fn: HandlerFunc(func(resp *Response, req *Request) {
				resp.SetStatus(200)
			}),
			expected: expectation{
				status:      200,
				body:        `{}` + "\n",
				contentType: "application/json",
			},
		},
		{
			body: `{"name":"Eric Holmes"}`,
			fn: HandlerFunc(func(resp *Response, req *Request) {
				var f UserCreateForm
				req.Decode(&f)

				resp.SetStatus(200)
				resp.Present(f)
			}),
			expected: expectation{
				status:      200,
				body:        `{"name":"Eric Holmes"}` + "\n",
				contentType: "application/json",
			},
		},
	}

	for _, test := range tests {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "", bytes.NewReader([]byte(test.body)))

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
