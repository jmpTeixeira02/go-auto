package sqlite

import (
	"fmt"
	"go-auto/data"
	dataService "go-auto/service/data"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Car struct {
	Car          string
	Price        int
	Mileage      int
	Fuel         string
	Year         int
	Hp           int
	Displacement int
	Link         string `gorm:"primaryKey"`
}

type SQLService struct {
	db *gorm.DB
}

func New(p string) (dataService.DataService, error) {
	db, err := gorm.Open(sqlite.Open(p), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the db %w", err)
	}

	err = db.AutoMigrate(&data.Car{})
	if err != nil {
		return nil, fmt.Errorf("Failed to migrate db schema %w", err)
	}

	return SQLService{
		db: db,
	}, nil
}

func (s SQLService) FindNewCars(cars []data.Car) ([]data.Car, error) {
	links := make([]string, len(cars))
	for i := range cars {
		links[i] = cars[i].Link
	}
	var table []Car

	// Need to check if DB is empty first, since the NOT IN
	// operation returns no values if the DB is empty
	var count int64
	s.db.Model(&table).Count(&count)
	if count == 0 {
		return cars, nil
	}

	var res []data.Car
	s.db.Where("link NOT IN ?", links).Find(&table).Scan(&res)
	return res, s.db.Error
}

func (s SQLService) AddCar(car data.Car) (data.Car, error) {
	model := carToCarDataModel(car)
	s.db.Create(&model)
	return carDataModelToCar(model), s.db.Error
}

func carToCarDataModel(car data.Car) Car {
	return Car{
		Car:          car.Car,
		Price:        car.Price,
		Mileage:      car.Mileage,
		Fuel:         car.Fuel,
		Year:         car.Year,
		Hp:           car.Hp,
		Displacement: car.Displacement,
		Link:         car.Link,
	}
}

func carDataModelToCar(car Car) data.Car {
	return data.Car{
		Car:          car.Car,
		Price:        car.Price,
		Mileage:      car.Mileage,
		Fuel:         car.Fuel,
		Year:         car.Year,
		Hp:           car.Hp,
		Displacement: car.Displacement,
		Link:         car.Link,
	}
}
