package data

import (
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

func CarScrapperToCar(carScrape scrapper.CarScrape) Car {
	car := Car{
		Car:     carScrape.Model,
		Price:   stringToInt(carScrape.Price),
		Mileage: stringToInt(strings.Split(carScrape.Mileage, "k")[0]),
		Fuel:    carScrape.Fuel,
		Year:    stringToInt(carScrape.Year),
		Link:    carScrape.Link,
	}
	err := UpdateCarPower(&car, carScrape.Power)
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
