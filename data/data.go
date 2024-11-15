package data

import (
	"fmt"
	"go-auto/scrapper"
	"strconv"
	"strings"
)

type Car struct {
	Car          string
	Price        int
	Mileage      int
	Fuel         string
	Year         int
	Hp           int
	Displacement int
	Link         string
}

func CarToString(car Car) string {
	return fmt.Sprintf("%s | %d € | %d KM | %s | %d | %d HP | %d cm3\n%s\n", car.Car, car.Price, car.Mileage, car.Fuel, car.Year, car.Hp, car.Displacement, car.Link)
}

func separateHpDisplace(str string) []string {
	return strings.Split(str, "•")
}

func stringToInt(str string) int {
	res, err := strconv.Atoi(strings.ReplaceAll(str, " ", ""))
	if err != nil {
		panic(err)
	}
	return res
}

func CarScrapperToCar(scrapped scrapper.Car) Car {
	car := Car{
		Car:     scrapped.Model,
		Price:   stringToInt(scrapped.Price),
		Mileage: stringToInt(strings.Split(scrapped.Mileage, "k")[0]),
		Fuel:    scrapped.Fuel,
		Year:    stringToInt(scrapped.Year),
		Link:    scrapped.Link,
	}
	err := UpdateCarPower(&car, scrapped.Power)
	if err != nil {
		panic(err)
	}
	return car
}

func UpdateCarPower(car *Car, raw string) error {
	res := make([]int, 2)
	for i, str := range separateHpDisplace(raw) {
		res[i] = stringToInt(strings.Split(str, "c")[0])
	}
	car.Displacement = res[0]
	car.Hp = res[1]
	return nil
}
