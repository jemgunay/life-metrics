package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	functions "github.com/jemgunay/life-metrics"
)

func main() {
	// define handlers
	http.HandleFunc("/daylog/data/daylog", enableCORS(functions.DayLogHandler))
	http.HandleFunc("/health", healthHandler)

	log.Printf("HTTP server starting on port %d", functions.Conf.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(functions.Conf.Port), nil)
	log.Printf("HTTP server shut down: %s", err)
}

var startTimestamp = time.Now().UTC().Format(time.RFC3339)

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Life-Metrics-Start-Time", startTimestamp)
	w.WriteHeader(http.StatusOK)
}

// enableCORS enables CORS for handlers that it wraps.
func enableCORS(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}
