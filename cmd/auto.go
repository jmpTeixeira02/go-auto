package main

import (
	"go-auto/config"
	"go-auto/notifier"
	"go-auto/service"
	"go-auto/service/data/sqlite"
)

func main() {
	conf, err := config.GetConf()
	if err != nil {
		panic(err)
	}
	notifier, err := notifier.NewNotifier(conf.Notifier)
	if err != nil {
		panic(err)
	}
	data, err := sqlite.New("cars.db")
	if err != nil {
		panic(err)
	}
	s := service.New(notifier, data)
	err = s.GetCars(conf.Url)
	if err != nil {
		panic(err)
	}
}
