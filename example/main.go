package main

import (
	"net/http"

	"github.com/ejholmes/buble"
	"github.com/gorilla/mux"
)

// Car is our domain model.
type Car struct {
	Make, Model string
}

var cars = []Car{
	{"Mini", "Cooper"},
	{"Honda", "Civic"},
	{"Toyota", "Celica"},
}

func main() {
	r := mux.NewRouter()
	r.Handle("/cars", &buble.Handler{
		HandlerFunc: CarsInfo,
	}).Methods("GET")
	r.Handle("/cars", &buble.Handler{
		HandlerFunc: CarsCreate,
	}).Methods("POST")

	http.ListenAndServe(":3000", r)
}

// CarsCreate adds a new car to the list of cars.
func CarsCreate(resp *buble.Response, req *buble.Request) {
	var c Car
	req.Decode(&c)

	cars = append(cars, c)

	resp.SetStatus(200)
	resp.Present(c)
}

// CarsInfo presents all of the cars.
func CarsInfo(resp *buble.Response, req *buble.Request) {
	resp.SetStatus(200)
	resp.Present(cars)
}
