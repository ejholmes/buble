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
func CarsCreate(w buble.ResponseWriter, r *buble.Request) {
	var c Car
	r.Decode(&c)
	cars = append(cars, c)

	w.WriteHeader(200)
	w.Encode(c)
}

// CarsInfo presents all of the cars.
func CarsInfo(w buble.ResponseWriter, r *buble.Request) {
	w.WriteHeader(200)
	w.Encode(cars)
}
