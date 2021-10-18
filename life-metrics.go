package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jemgunay/life-metrics/api"
	"github.com/jemgunay/life-metrics/config"
	"github.com/jemgunay/life-metrics/influx"
	"github.com/jemgunay/life-metrics/sources"
	"github.com/jemgunay/life-metrics/sources/monzo"
)

func main() {
	pollInterval := flag.Duration("poll_interval", time.Minute*10, "how often to poll the sources")
	flag.Parse()

	conf := config.New()

	// influx storage
	influxRequester := influx.New(conf.InfluxHost, conf.InfluxToken)

	// configure data sources
	monzoSource := monzo.New()
	dataSources := []sources.Source{
		monzoSource,
	}
	// configure the API
	api := api.New(influxRequester)

	// poll and scrape data from sources at fixed interval
	endTime := time.Now().UTC()
	startTime := endTime.Add(-*pollInterval)
	pollChan := make(chan struct{}, 1)
	go func() {
		ticker := time.NewTicker(*pollInterval)
		for {
			select {
			case <-pollChan:
			case <-ticker.C:
			}

			// perform collection
			wg := &sync.WaitGroup{}
			wg.Add(len(dataSources))
			for _, source := range dataSources {
				go func(source sources.Source) {
					defer wg.Done()

					// perform source collection
					log.Printf("collecting from source: %s", source.Name())
					results, err := source.Collect(startTime, endTime)
					if err != nil {
						log.Printf("source collection failed: %s: %s", source.Name(), err)
						return
					}

					// no new data to store so skip
					if len(results) == 0 {
						return
					}

					// write collected source data to influx
					if err := influxRequester.Write(source.Name(), results...); err != nil {
						log.Printf("writing source data to influx failed: %s: %s", source.Name(), err)
					}
				}(source)
			}
		}
	}()

	// define handlers
	http.HandleFunc("/api/data", enableCORS(api.Handler))
	http.HandleFunc("/api/data/collect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		pollChan <- struct{}{}
	})
	http.HandleFunc("/api/auth/monzo", monzoSource.StartOauth)
	http.HandleFunc("/api/auth/monzo/callback", monzoSource.CompleteOauth)
	http.HandleFunc("/health", healthHandler)

	log.Printf("HTTP server starting on port %d", conf.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port), nil)
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
