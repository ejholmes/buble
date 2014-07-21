package main

import (
	"net/http"

	"github.com/ejholmes/buble"
	"github.com/gorilla/mux"
)

var cars = []Car{
	{"Mini", "Cooper"},
	{"Honda", "Civic"},
	{"Toyota", "Celica"},
}

// Car is our domain model.
type Car struct {
	Make, Model string
}

// ResponseWriter wraps buble.ResponseWriter.
type ResponseWriter interface {
	buble.ResponseWriter
}

// Response wraps buble.Response.
type Response struct {
	buble.ResponseWriter
}

// Request wraps buble.Request.
type Request struct {
	*buble.Request
}

// HandlerFunc defines our handler function signature.
type HandlerFunc func(ResponseWriter, *Request)

// API implements the http.Handler interface for serving requests.
type API struct {
	router *mux.Router
}

// NewAPI returns a new API.
func NewAPI() *API {
	return &API{router: mux.NewRouter()}
}

// Handle creates a new route that will be handled using fn.
func (a *API) Handle(method, path string, fn HandlerFunc) {
	h := &buble.Handler{
		HandlerFunc: func(w buble.ResponseWriter, r *buble.Request) {
			fn(&Response{ResponseWriter: w}, &Request{Request: r})
		},
	}
	a.router.Handle(path, h).Methods(method)
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func main() {
	a := NewAPI()

	a.Handle("GET", "/cars", CarsInfo)
	a.Handle("POST", "/cars", CarsCreate)

	http.ListenAndServe(":3000", a)
}

// CarsCreate adds a new car to the list of cars.
func CarsCreate(w ResponseWriter, r *Request) {
	var c Car
	r.Decode(&c)
	cars = append(cars, c)

	w.WriteHeader(200)
	w.Encode(c)
}

// CarsInfo presents all of the cars.
func CarsInfo(w ResponseWriter, r *Request) {
	w.WriteHeader(200)
	w.Encode(cars)
}
