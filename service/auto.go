package service

import (
	"go-auto/data"
	"go-auto/notifier"
	"go-auto/scrapper"
)

func GetCars(url string, notifier notifier.Notifier) error {
	s := scrapper.New()

	var carsScrape []scrapper.Car
	err := s.Scrape(
		url,
		&carsScrape,
		scrapper.GetCarModel,
		scrapper.GetCarPower,
		scrapper.GetCarDetails,
		scrapper.GetCarPrice,
	)
	if err != nil {
		return err
	}
	cars := make([]data.Car, len(carsScrape))
	for i, car := range carsScrape {
		cars[i] = data.CarScrapperToCar(car)
	}

	// TODO -> Save cars and only notify on new ones

	for _, car := range cars {
		go notifier.SendMessage(data.CarToString(car))
		// Add to db if new car was sucessufly saved
	}

	return nil
}
