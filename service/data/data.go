package data

import "go-auto/data"

type DataService interface {
	FindNewCars(cars []data.Car) ([]data.Car, error)
	AddCar(data.Car) (data.Car, error)
}
