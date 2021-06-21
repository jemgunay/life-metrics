package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jemgunay/life-metrics/api"
	"github.com/jemgunay/life-metrics/influx"
	"github.com/jemgunay/life-metrics/sources"
	"github.com/jemgunay/life-metrics/sources/monzo"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("no ENV PORT defined - defaulting to port %s", port)
	}
	influxToken := os.Getenv("INFLUX_TOKEN")
	if influxToken == "" {
		influxToken = "zvdEmxrLzHunj6n2PxgmMsqubwoTjfZEBJFDKMSZIqoBZ2pe09_W9-JY9TYxQj3_oP2q8pb_HBLO3_QMufSNLw=="
		log.Printf("no ENV INFLUX_TOKEN defined - defaulting to %s", port)
	}
	influxHost := os.Getenv("INFLUX_HOST")
	if influxHost == "" {
		influxHost = "http://localhost:8086"
		log.Printf("no ENV INFLUX_HOST defined - defaulting to %s", port)
	}
	pollInterval := flag.Duration("poll_interval", time.Minute*10, "how often to poll the sources")
	flag.Parse()

	// configure the API
	api := api.New()
	// configure data sources
	dataSources := []sources.Source{
		 monzo.New(),
	}

	// influx storage
	influxRequester := influx.New(influxHost, influxToken)

	// poll and scrape data from sources at fixed interval
	startTime := time.Now().UTC().Add(-*pollInterval)
	endTime := time.Now().UTC()
	pollChan := make(chan struct{}, 1)
	go func() {
		ticker := time.NewTicker(*pollInterval)
		for {
			select {
			case logPayload := <-api.Updates():
				influxRequester.Write("day_log", logPayload)
			case <-pollChan:
			case <-ticker.C:
			}

			// perform collection
			wg := &sync.WaitGroup{}
			wg.Add(len(dataSources))
			for _, source := range dataSources {
				go func() {
					defer wg.Done()

					// perform source collection
					log.Printf("collecting from source: %s", source.Name())
					results, err := source.Collect(startTime, endTime)
					if err != nil {
						log.Printf("source collection failed: %s", err)
						return
					}

					influxRequester.Write(source.Name(), results...)
				}()
			}
		}
	}()

	// define handlers
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/flush", func(w http.ResponseWriter, r *http.Request) {
		pollChan <- struct{}{}
	})
	http.HandleFunc("/api", enableCORS(api.Handler))

	log.Printf("HTTP server starting on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Printf("HTTP server shut down: %s", err)
}

func enableCORS(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

var startTimestamp = time.Now().UTC().String()

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Life-Metrics-Start-Time", startTimestamp)
	w.WriteHeader(http.StatusOK)
}
