package main

import (
	"encoding/json"
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

// TODO: start API here for local development of functions
func main() {
	conf := config.New()

	// influx storage
	influxRequester := influx.New(conf.Influx)

	// configure data sources
	monzoSource := monzo.New(conf, influxRequester)
	p := poller{
		sources: []sources.Source{
			monzoSource,
		},
		scrapeChan: make(chan collectRequest, 1),
	}

	// start collection poller
	go p.start(influxRequester)

	// define handlers
	apiHandler := api.New(influxRequester).Handler
	http.HandleFunc("/api/data/daylog", enableCORS(apiHandler))
	http.HandleFunc("/api/data/collect", enableCORS(p.collectHandler))
	http.HandleFunc("/api/data/sources", enableCORS(p.sourcesHandler))
	http.HandleFunc("/api/auth/monzo", monzoSource.AuthenticateHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("HTTP server starting on port %d", conf.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port), nil)
	log.Printf("HTTP server shut down: %s", err)
}

// collectRequest specifies collection details.
type collectRequest struct {
	reset bool
}

// poller serialises access to source operations.
type poller struct {
	sources    []sources.Source
	scrapeChan chan collectRequest
}

// start polls for scrape requests and performs collections for each source.
func (p poller) start(influxRequester influx.Requester) {
	for req := range p.scrapeChan {
		endTime := time.Now().UTC()

		// perform collection for each source
		for _, source := range p.sources {
			var startTime time.Time
			if req.reset {
				startTime = time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)

			} else {
				var err error
				startTime, err = influxRequester.LastTimestampByMeasurement(source.Name())
				if err != nil {
					log.Printf("failed to get last timestamp for source %s: %s", source.Name(), err)
					continue
				}
				// add a second to ensure we don't recollect the last record
				startTime = startTime.Add(time.Second)
			}

			source.Collect(sources.NewPeriod(startTime, endTime))
		}
	}
}

func (p poller) collectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req := collectRequest{
		reset: r.URL.Query().Get("reset") == "true",
	}

	select {
	case p.scrapeChan <- req:
		w.WriteHeader(http.StatusAccepted)
	default:
		w.WriteHeader(http.StatusTooManyRequests)
	}
}

func (p poller) sourcesHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]sources.StateSet, len(p.sources))
	for _, source := range p.sources {
		resp[source.Name()] = source.State()
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Printf("failed to JSON marshal source state: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
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
