package main

import (
	"fmt"
	"go-auto/config"
	"go-auto/notifier"
	"go-auto/service"
	"go-auto/service/data/sqlite"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	data, err := sqlite.New(conf.Data.Config.Address)
	if err != nil {
		panic(err)
	}
	s := service.New(notifier, data)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(time.Duration(conf.Refresh * int(time.Minute)))
	defer ticker.Stop()

	err = s.GetCars(conf.Url)
	if err != nil {
		panic(err)
	}
	go func() {
		for range ticker.C {
			err = s.GetCars(conf.Url)
			if err != nil {
				panic(err)
			}
		}
	}()
	<-sigChan
	fmt.Println("Closing the application")
}
