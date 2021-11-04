package config

import (
	"log"
	"os"
	"strconv"
)

// Config is the service config.
type Config struct {
	Port              int
	InfluxHost        string
	InfluxToken       string
	MonzoClientID     string
	MonzoClientSecret string
}

// New initialises a Config from environment variables.
func New() Config {
	// attempt to get config environment vars, or default them
	return Config{
		Port:              getEnvVarInt("PORT", 8080),
		InfluxHost:        getEnvVar("INFLUX_HOST", "http://localhost:8086"),
		InfluxToken:       getEnvVar("INFLUX_TOKEN", ""),
		MonzoClientID:     getEnvVar("MONZO_CLIENT_ID", ""),
		MonzoClientSecret: getEnvVar("MONZO_CLIENT_SECRET", ""),
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
