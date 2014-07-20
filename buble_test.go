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
				w.WriteHeader(200)
				w.Encode(&User{ID: 1})
			}),
			expected: expectation{
				status:      200,
				body:        `{"id":1}` + "\n",
				contentType: "application/json",
			},
		},
		{
			method: "POST",
			body:   `{"name":"Eric Holmes"}`,
			fn: HandlerFunc(func(w ResponseWriter, r *Request) {
				var f UserCreateForm
				r.Decode(&f)

				w.WriteHeader(200)
				w.Encode(f)
			}),
			expected: expectation{
				status:      200,
				body:        `{"name":"Eric Holmes"}` + "\n",
				contentType: "application/json",
			},
		},
	}

	for i, test := range tests {
		method := test.method
		if method == "" {
			method = "GET"
		}

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(method, "", bytes.NewReader([]byte(test.body)))

		h := &Handler{HandlerFunc: test.fn}
		h.ServeHTTP(resp, req)

		if resp.Code != test.expected.status {
			t.Errorf("Status %v: Want %v, Got %v", i, test.expected.status, resp.Code)
		}

		if resp.Body.String() != test.expected.body {
			t.Errorf("Body %v: Want %v, Got %v", i, test.expected.body, resp.Body.String())
		}

		contentType := resp.HeaderMap.Get("Content-Type")
		if contentType != test.expected.contentType {
			t.Errorf("Content-Type %v: Want %v, Got %v", i, test.expected.contentType, contentType)
		}
	}
}
