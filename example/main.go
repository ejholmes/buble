package main

import (
	"net/http"

	"github.com/ejholmes/buble"
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
	http.Handle("/cars", &buble.Handler{
		HandlerFunc: GetCars,
	})

	http.ListenAndServe(":3000", nil)
}

// GetCars presents all of the cars.
func GetCars(resp *buble.Response, req *buble.Request) {
	resp.SetStatus(200)
	resp.Present(cars)
}
