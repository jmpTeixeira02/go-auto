package sqlite

import (
	"go-auto/data"
	dataService "go-auto/service/data"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initTest() (dataService.DataService, func()) {
	testFile := "test"
	service, err := New(testFile)
	if err != nil {
		panic(err)
	}
	return service, func() {
		err = os.Remove(testFile)
		if err != nil {
			panic(err)
		}
	}
}

func TestAddCar(t *testing.T) {
	service, tearDown := initTest()
	defer tearDown()
	input := data.Car{
		Car:          "",
		Price:        0,
		Mileage:      0,
		Fuel:         "",
		Year:         0,
		Hp:           0,
		Displacement: 0,
		Link:         "",
	}
	car, err := service.AddCar(input)
	assert.Nil(t, err)
	assert.EqualValues(t, input, car)
}

func TestFindNewCars(t *testing.T) {
	service, tearDown := initTest()
	defer tearDown()

	input := data.Car{
		Car:          "",
		Price:        0,
		Mileage:      0,
		Fuel:         "",
		Year:         0,
		Hp:           0,
		Displacement: 0,
		Link:         "test",
	}

	_, err := service.AddCar(input)
	assert.Nil(t, err)

	cars, err := service.FindNewCars([]data.Car{{Link: ""}})
	assert.Nil(t, err)
	assert.Len(t, cars, 1)

	cars, err = service.FindNewCars([]data.Car{{Link: "test"}})
	assert.Nil(t, err)
	assert.Len(t, cars, 0)
}
