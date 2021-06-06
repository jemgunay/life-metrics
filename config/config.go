package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Monzo `json:"monzo"`
}

type Monzo struct {
	Host string
	Key string
}

func New(configFile string) (Config, error) {
	conf := Config{}
	b, err := ioutil.ReadFile("config/" + configFile)
	if err != nil {
		return conf, err
	}

	if err := json.Unmarshal(b, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}