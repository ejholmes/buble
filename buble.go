package buble

import (
	"encoding/json"
	"net/http"
)

// Form is an interface for validating incoming POST/PUT/PATCH data.
type Form interface {
	Validate() error
}

// Resource is an interface that provides a means for transforming
// a domain model into an API resource suitable for encoding.
type Resource interface {
	Present() interface{}
}

// Request represents the http.Request and passes along the decoded form
// associated with the Handler.
type Request struct {
	Form Form
}

// Response provides a means for building an http response.
type Response struct {
	Resource Resource
	Status   int
}

// HandlerFunc is a function signature for handling a request.
type HandlerFunc func(*Response, *Request)

// Handler represents an API endpoint and implements the http.Handler
// interface for serving a request.
type Handler struct {
	Formatter   Formatter
	Form        Form
	Resource    Resource
	HandlerFunc HandlerFunc
}

// ServeHTTP implements the http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.HandlerFunc == nil {
		panic("no HandlerFunc provided")
	}

	resp := &Response{}
	req := &Request{}

	h.HandlerFunc(resp, req)

	w.WriteHeader(resp.Status)
}

// Formatter is an interface for encoding/decoding requests/responses.
type Formatter interface {
	// Decode takes an http request, decodes the body into Form.
	Decode(*http.Request, Form)

	// Encode encodes the Resource into the response.
	Encode(Resource, http.ResponseWriter)
}

// JSONFormatter is an implementation of the Formatter interface for
// encoding/decoding JSON.
type JSONFormatter struct {
}

// Decode decodes the request body into form.
func (fmtr *JSONFormatter) Decode(r *http.Request, f Form) {
	json.NewDecoder(r.Body).Decode(f)
}

// Encode encodes the Resource into the response and sets the
// Content-Type header to "application/json".
func (fmtr *JSONFormatter) Encode(r Resource, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
}
