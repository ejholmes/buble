package buble

import (
	"encoding/json"
	"net/http"
)

// DefaultFormatter is the default formatter to use in a Handler
// if no formatter is specified.
var DefaultFormatter = &JSONFormatter{}

// Request wraps http.Request.
type Request struct {
	*http.Request
	Decoder Decoder
}

// Decode decodes the request body into v.
func (r *Request) Decode(v interface{}) error {
	return r.Decoder.Decode(r.Request, v)
}

// ResponseWriter is an interface that wraps http.ResponseWriter with some convenience.
type ResponseWriter interface {
	http.ResponseWriter
	SetStatus(int)
	Present(interface{})
}

// Response is an implementation of the ResponseWriter inteface.
type Response struct {
	http.ResponseWriter
	Encoder Encoder

	Resource interface{}
	Status   int
}

// SetStatus sets the status code for the response.
func (r *Response) SetStatus(status int) {
	r.Status = status
}

// Present sets the resource.
func (r *Response) Present(v interface{}) {
	r.Resource = v
}

// HandlerFunc is a function signature for handling a request.
type HandlerFunc func(ResponseWriter, *Request)

// Handler represents an API endpoint and implements the http.Handler
// interface for serving a request.
type Handler struct {
	// Formatter is the formatter to use when encoding/decoding the request/response.
	// The zero value is the DefaultFormatter.
	Formatter Formatter

	// HandlerFunc is the HandlerFunc to call when a request is handled.
	HandlerFunc HandlerFunc
}

// ServeHTTP implements the http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := &Request{Request: r, Decoder: h.formatter()}
	resp := &Response{ResponseWriter: w, Encoder: h.formatter()}

	h.handlerFunc()(resp, req)

	resp.WriteHeader(resp.Status)
	resp.Encoder.Encode(resp.Resource, resp.ResponseWriter)
}

func (h *Handler) formatter() Formatter {
	if h.Formatter == nil {
		h.Formatter = DefaultFormatter
	}
	return h.Formatter
}

func (h *Handler) handlerFunc() HandlerFunc {
	if h.HandlerFunc == nil {
		panic("no HandlerFunc provided")
	}
	return h.HandlerFunc
}

// Encoder is an interface for encoding a value into an http response.
type Encoder interface {
	Encode(interface{}, http.ResponseWriter) error
}

// Decoder is an interface for decoding the request body into an interface.
type Decoder interface {
	Decode(*http.Request, interface{}) error
}

// Formatter is an interface for encoding/decoding requests/responses.
type Formatter interface {
	Encoder
	Decoder
}

// JSONFormatter is an implementation of the Formatter interface for
// encoding/decoding JSON.
type JSONFormatter struct {
}

// Decode decodes the request body into form.
func (fmtr *JSONFormatter) Decode(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(v)
}

// Encode encodes the Resource into the response and sets the
// Content-Type header to "application/json".
func (fmtr *JSONFormatter) Encode(v interface{}, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")

	// If we set the Content-Type to application/json, we should always respond with valid json.
	if v == nil {
		v = map[string]interface{}{}
	}

	return json.NewEncoder(w).Encode(v)
}
