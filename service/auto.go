package service

import (
	"go-auto/data"
	"go-auto/scrapper"
)

func GetCars(url string) []data.Car {
	s := scrapper.New()

	var carsScrape []scrapper.CarScrape
	s.Scrape(
		url,
		&carsScrape,
		scrapper.GetCarModel,
		scrapper.GetCarPower,
		scrapper.GetCarDetails,
		scrapper.GetCarPrice,
	)
	cars := make([]data.Car, len(carsScrape))
	for i, car := range carsScrape {
		cars[i] = data.CarScrapperToCar(car)
	}
	return cars
}
