package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jemgunay/life-metrics/config"
	"github.com/jemgunay/life-metrics/influx"
	"github.com/jemgunay/life-metrics/sources"
	"github.com/jemgunay/life-metrics/sources/monzo"
)

func main() {
	port := flag.Int("port", 8080, "HTTP server port")
	pollInterval := flag.Duration("poll_interval", time.Minute*10, "how often to poll the sources")
	flag.Parse()

	conf, err := config.New()
	if err != nil {
		log.Printf("failed to read config: %s", err)
		return
	}

	// configure data sources
	monzoSource := monzo.New(conf.Monzo)

	dataSources := map[string]sources.Source{
		"monzo": monzoSource,
	}

	// influx storage
	influxRequester := influx.New()

	// poll and scrape data from sources at fixed interval
	startTime := time.Now().UTC().Add(-*pollInterval)
	endTime := time.Now().UTC()
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
			for sourceName, sourceType := range dataSources {
				go func() {
					defer wg.Done()

					// perform source collection
					log.Printf("collecting from source: %s", sourceName)
					results, err := sourceType.Collect(startTime, endTime)
					if err != nil {
						log.Printf("source collection failed: %s", err)
						return
					}

					influxRequester.Write(results)
				}()
			}
		}
	}()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/flush", func(w http.ResponseWriter, r *http.Request) {
		pollChan <- struct{}{}
	})
	err = http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	log.Printf("HTTP server shut down: %s", err)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
