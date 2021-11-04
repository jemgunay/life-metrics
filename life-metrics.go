package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jemgunay/life-metrics/api"
	"github.com/jemgunay/life-metrics/config"
	"github.com/jemgunay/life-metrics/influx"
	"github.com/jemgunay/life-metrics/sources"
	"github.com/jemgunay/life-metrics/sources/monzo"
)

func main() {
	conf := config.New()

	// influx storage
	influxRequester := influx.New(conf.InfluxHost, conf.InfluxToken)

	// configure data sources
	monzoSource := monzo.New(conf, influxRequester)
	dataSources := []sources.Source{
		monzoSource,
	}

	// poll and scrape data from sources at fixed interval
	pollChan := make(chan bool, 1)
	go func() {
		for reset := range pollChan {
			// TODO: allow specify manual range & source name in request for force range collections
			endTime := time.Now().UTC()

			// perform source collection
			for _, source := range dataSources {
				var startTime time.Time
				if reset {
					startTime = time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)
				} else {
					var err error
					startTime, err = influxRequester.LastTimestampByMeasurement(source.Name())
					if err != nil {
						log.Printf("failed to get last timestamp for source %s: %s", source.Name(), err)
						continue
					}
				}

				source.Collect(sources.NewPeriod(startTime, endTime))
			}
		}
	}()

	// define handlers
	apiHandler := api.New(influxRequester).Handler
	http.HandleFunc("/api/data", enableCORS(apiHandler))
	http.HandleFunc("/api/data/collect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		reset := r.URL.Query().Get("reset") == "true"
		// TODO: url query for start/end/source name
		pollChan <- reset
	})
	http.HandleFunc("/api/auth/monzo", monzoSource.AuthenticateHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("HTTP server starting on port %d", conf.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port), nil)
	log.Printf("HTTP server shut down: %s", err)
}

// enableCORS enables CORS for handlers that it wraps.
func enableCORS(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

var startTimestamp = time.Now().UTC().Format(time.RFC3339)

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Life-Metrics-Start-Time", startTimestamp)
	w.WriteHeader(http.StatusOK)
}
