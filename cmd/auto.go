package main

import (
	"fmt"
	"go-auto/service"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Url      string `yaml:"url"`
	BotToken string `yaml:"botToken"`
	User     string `yaml:"user"`
}

func (c *Config) getConf() *Config {
	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		panic(fmt.Errorf("error opening the yaml %w", err))
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic(fmt.Errorf("error unmarshalling the yaml %w", err))
	}

	return c
}

func main() {
	var c Config
	c.getConf()
	if c.getConf().Url == "" {
		panic("You must set the url in the config file!")
	}
	if c.getConf().Url == "" {
		panic("You must set the discord token in the config file!")
	}
	if c.getConf().Url == "" {
		panic("You must set the user in the config file!")
	}
	cars := service.GetCars(c.getConf().Url)
	fmt.Println(cars)
}
