package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Url      string   `yaml:"url"`
	Refresh  int      `yaml:"refresh_min"`
	Notifier Notifier `yaml:"notifier"`
	Data     Data     `yaml:"data"`
}

type Notifier struct {
	Service string         `yaml:"service"`
	Config  notifierConfig `yaml:"config"`
}

type Data struct {
	Service string     `yaml:"service"`
	Config  dataConfig `yaml:"config"`
}

type dataConfig struct {
	Address string `yaml:"address"`
}

type notifierConfig struct {
	Token    string `yaml:"token"`
	Receiver string `yaml:"receiver"`
}

func GetConf() (*Config, error) {
	var conf Config
	yamlFile, err := os.ReadFile("config/config.yml")
	if err != nil {
		return nil, fmt.Errorf("error opening the yaml %w", err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling the yaml %w", err)
	}

	return &conf, err
}
