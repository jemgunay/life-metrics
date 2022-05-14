package main

import (
	"log"
	"net/http"
	"strconv"

	functions "github.com/jemgunay/life-metrics"
)

func main() {
	http.HandleFunc("/daylog/data/daylog", functions.DayLogHandler)

	log.Printf("HTTP server starting on port %d", functions.Conf.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(functions.Conf.Port), nil)
	log.Printf("HTTP server shut down: %s", err)
}