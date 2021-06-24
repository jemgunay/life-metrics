package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	InfluxHost  string
	InfluxToken string
}

func New() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("no ENV PORT defined - defaulting to port %s", port)
	}

	influxHost := os.Getenv("INFLUX_HOST")
	if influxHost == "" {
		influxHost = "http://localhost:8086"
		log.Printf("no ENV INFLUX_HOST defined - defaulting to %s", influxHost)
	} else {
		log.Printf("ENV INFLUX_HOST found")
	}

	influxToken := os.Getenv("INFLUX_TOKEN")
	if influxToken == "" {
		influxToken = "***REMOVED***"
		log.Printf("no ENV INFLUX_TOKEN defined - defaulting to %s", influxToken)
	} else {
		log.Printf("ENV INFLUX_TOKEN found")
	}

	return Config{
		Port:        port,
		InfluxHost:  influxHost,
		InfluxToken: influxToken,
	}
}
