package data

import (
	"go-auto/scrapper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCarPower(t *testing.T) {
	// Arrange
	var actual Car
	expected := Car{
		Hp:           100,
		Displacement: 1342,
	}

	// Act
	err := UpdateCarPower(&actual, "1 342 cm3 • 100 cv")
	if err != nil {
		panic(err)
	}

	// Assert
	assert.EqualValues(t, expected, actual)
}

func TestCarScrapperToCar(t *testing.T) {
	// Arrange
	expected := Car{
		Car:          "Toyota Yaris",
		Price:        10000,
		Mileage:      89123,
		Fuel:         "Gasolina",
		Year:         2008,
		Hp:           100,
		Displacement: 1342,
		Link:         "test",
	}
	scrapped := scrapper.Car{
		Model:   "Toyota Yaris",
		Price:   "10 000",
		Mileage: "89 123",
		Fuel:    "Gasolina",
		Year:    "2008",
		Power:   "1 342 cm3 • 100 cv",
		Link:    "test",
	}

	// Act
	car := CarScrapperToCar(scrapped)

	// Assert
	assert.EqualValues(t, expected, car)
}
