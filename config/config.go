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

func New() (Config, error) {
	conf := Config{}
	b, err := ioutil.ReadFile("config/LOCAL.JSON")
	if err != nil {
		return conf, err
	}

	if err := json.Unmarshal(b, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}