package config

import (
	"log"
	"os"
	"strconv"
)

// Config is the service config.
type Config struct {
	Port   int
	Influx Influx
	Monzo  Monzo
}

// Influx contains the InfluxDB config.
type Influx struct {
	Host  string
	Token string
	Org   string
}

// Monzo contains the Monzo source config.
type Monzo struct {
	ClientID     string
	ClientSecret string
	RedirectHost string
}

// New initialises a Config from environment variables.
func New() Config {
	// attempt to get config environment vars, or default them
	return Config{
		Port: getEnvVarInt("PORT", 8080),
		Influx: Influx{
			Host:  getEnvVar("INFLUX_HOST", "http://localhost:8086"),
			Token: getEnvVar("INFLUX_TOKEN", ""),
			Org:   getEnvVar("INFLUX_ORG", ""),
		},
		Monzo: Monzo{
			ClientID:     getEnvVar("MONZO_CLIENT_ID", ""),
			ClientSecret: getEnvVar("MONZO_CLIENT_SECRET", ""),
			RedirectHost: getEnvVar("MONZO_REDIRECT_HOST", "http://localhost:8080"),
		},
	}
}

// getEnvVar gets a string environment variable or defaults it if unset.
func getEnvVar(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("no %s environment var defined - defaulting to %s", key, defaultValue)
		return defaultValue
	}

	log.Printf("%s environment var found", key)
	return val
}

// getEnvVarInt gets an integer environment variable or defaults it if unset.
func getEnvVarInt(key string, defaultValue int) int {
	varStr := getEnvVar(key, strconv.Itoa(defaultValue))
	varInt, _ := strconv.Atoi(varStr)
	return varInt
}
