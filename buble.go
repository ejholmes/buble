package buble

import "net/http"

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

// ServeHTTP implements the http.Handler interface.
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	req := &Request{}

	h(resp, req)
}

// Handler represents an API endpoint and implements the http.Handler
// interface for serving a request.
type Handler struct {
	Form        Form
	Resource    Resource
	HandlerFunc HandlerFunc
}

// ServeHTTP implements the http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.HandlerFunc == nil {
		panic("no HandlerFunc provided")
	}
	h.HandlerFunc.ServeHTTP(w, r)
}
