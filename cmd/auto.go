package main

import (
	"go-auto/config"
	"go-auto/notifier"
	"go-auto/service"
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
	service.GetCars(conf.Url, notifier)
}
