package service

import (
	"context"
	"errors"
	"go-auto/data"
	"go-auto/notifier"
	"go-auto/scrapper"
	dataService "go-auto/service/data"
	"log"
	"os"
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

	log *log.Logger
}

func New(notifier notifier.Notifier, dataService dataService.DataService) AutoService {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT)

	return AutoService{
		notifier: notifier,
		data:     dataService,
		ctx:      ctx,
		wg:       sync.WaitGroup{},
		log:      log.New(os.Stdout, "", 0),
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

func (a *AutoService) addNewCars(cars []data.Car) error {
	nCars := len(cars)
	a.wg.Add(nCars)
	a.errCh = make(chan error, nCars)
	for i := range cars {
		go a.addCar(cars[i])
	}
	a.wg.Wait()

	var err error
	for {
		select {
		case e := <-a.errCh:
			err = errors.Join(err, e)
		default:
			close(a.errCh)
			return err
		}
	}
}

func (a *AutoService) addCar(car data.Car) {
	defer a.wg.Done()
	err := a.notifier.SendMessage(data.CarToString(car))
	if err != nil {
		a.errCh <- err
		return
	}
	_, err = a.data.AddCar(car)
	if err != nil {
		a.errCh <- err
	}
}
