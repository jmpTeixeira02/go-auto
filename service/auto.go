package service

import (
	"go-auto/data"
	"go-auto/notifier"
	"go-auto/scrapper"
	dataService "go-auto/service/data"
)

type AutoService struct {
	notifier notifier.Notifier
	data     dataService.DataService
}

func New(notifier notifier.Notifier, dataService dataService.DataService) AutoService {
	return AutoService{
		notifier: notifier,
		data:     dataService,
	}
}

func (a *AutoService) GetCars(url string) error {
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

	// Process cars, give all and return only the non viewed ones
	new, err := a.data.FindNewCars(cars)
	if err != nil {
		return err
	}

	for i := range new {
		a.notifier.SendMessage(data.CarToString(new[i]))
		a.data.AddCar(new[i])
		// Add to db if new car was sucessufly saved
	}

	return nil
}
