package service

import (
	"context"
	"errors"
	"go-auto/data"
	"go-auto/notifier"
	"go-auto/scrapper"
	dataService "go-auto/service/data"
	"os/signal"
	"sync"
	"syscall"
)

type AutoService struct {
	notifier notifier.Notifier
	data     dataService.DataService

	ctx   context.Context
	errCh chan error
	wg    sync.WaitGroup
}

func New(notifier notifier.Notifier, dataService dataService.DataService) AutoService {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT)

	return AutoService{
		notifier: notifier,
		data:     dataService,
		ctx:      ctx,
		errCh:    make(chan error),
		wg:       sync.WaitGroup{},
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

	return a.addNewCars(new)
}

func (a *AutoService) addNewCars(new []data.Car) error {
	nNew := len(new)
	a.wg.Add(nNew)
	a.errCh = make(chan error)
	for i := range new {
		go a.addCar(new[i])
	}
	a.wg.Wait()

	var err error
	for e := range a.errCh {
		err = errors.Join(err, e)
	}
	return err
}

func (a *AutoService) addCar(car data.Car) {
	err := a.notifier.SendMessage(data.CarToString(car))
	if err != nil {
		a.errCh <- err
		return
	}
	_, err = a.data.AddCar(car)
	if err != nil {
		a.errCh <- err
		return
	}
	a.wg.Done()
}
