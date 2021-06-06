package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jemgunay/life-metrics/api"
	"github.com/jemgunay/life-metrics/config"
	"github.com/jemgunay/life-metrics/influx"
	"github.com/jemgunay/life-metrics/sources"
	"github.com/jemgunay/life-metrics/sources/monzo"
)

func main() {
	port := flag.Int("port", 8080, "HTTP server port")
	configFile := flag.String("config", "LOCAL.json", "config file")
	staticContentDir := flag.String("static_dir", "build/ui", "static UI assets dir")
	pollInterval := flag.Duration("poll_interval", time.Minute*10, "how often to poll the sources")
	flag.Parse()

	conf, err := config.New(*configFile)
	if err != nil {
		log.Printf("failed to read config: %s", err)
		return
	}

	// configure the API
	api := api.New()
	// configure data sources
	dataSources := []sources.Source{
		 monzo.New(conf.Monzo),
	}

	out, err := exec.Command("ls").Output()
	if err != nil {
		log.Printf("exec error: %s\n", err)
	} else {
		log.Println("Command Successfully Executed:", string(out))
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

	// define file server
	fileServerRoutes := []string{"/js/", "/css/", "/fonts/", "/favicon.ico"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, route := range fileServerRoutes {
			if strings.HasPrefix(r.URL.Path, route) {
				// serve static content
				http.ServeFile(w, r, *staticContentDir + r.URL.Path)
				return
			}
		}
		// serve Vue's index.html
		http.ServeFile(w, r, *staticContentDir + "/index.html")
	})

	log.Printf("HTTP server starting on port %d", *port)
	err = http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	log.Printf("HTTP server shut down: %s", err)
}

func enableCORS(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
