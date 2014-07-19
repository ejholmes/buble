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
	tests := []struct {
		fn       HandlerFunc
		expected struct {
			status      int
			body        string
			contentType string
		}
	}{
		{
			fn: HandlerFunc(func(resp *Response, req *Request) {
				resp.SetStatus(200)
				resp.Present(&User{ID: 1})
			}),
			expected: struct {
				status      int
				body        string
				contentType string
			}{
				status:      200,
				body:        `{"id":1}` + "\n",
				contentType: "application/json",
			},
		},
		{
			fn: HandlerFunc(func(resp *Response, req *Request) {
				resp.SetStatus(200)
			}),
			expected: struct {
				status      int
				body        string
				contentType string
			}{
				status:      200,
				body:        `{}` + "\n",
				contentType: "application/json",
			},
		},
	}

	for _, test := range tests {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "", nil)

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
