package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jemgunay/life-metrics/sources"
)

type request struct {
	Date    time.Time              `json:"date"`
	Metrics map[string]interface{} `json:"metrics"`
	Notes   string                 `json:"notes"`
}

var permittedMetrics = [...]string{
	"general_mood",
	"work_mood",
	"water_intake",
	"sleep_quality",
	"exercise",
	"meditation",
}

type API struct {
	updates chan sources.Result
}

func New() API {
	return API{
		updates: make(chan sources.Result, 1),
	}
}

func (a API) Handler(w http.ResponseWriter, r *http.Request) {
	// read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := request{}
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("failed to read body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := a.process(req); err != nil {
		log.Printf("failed to process request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (a API) process(req request) error {
	// validate metric fields
	if len(req.Metrics) != len(permittedMetrics) {
		return fmt.Errorf("expected %d metric fields", len(permittedMetrics))
	}

	var (
		tags       = make(map[string]string)
		scoreValue int
		scoreMax   int
	)
	for _, metricName := range permittedMetrics {
		rawMetric, ok := req.Metrics[metricName]
		if !ok {
			return fmt.Errorf("expected %s metric field", metricName)
		}

		switch metric := rawMetric.(type) {
		case int:
			if metric > 10 {
				return fmt.Errorf("integer metric %s field value is greater than 10", metricName)
			} else if metric < 0 {
				return fmt.Errorf("integer metric %s field value is less than 0", metricName)
			}
			tags[metricName] = strconv.Itoa(metric)
			scoreValue += metric
			scoreMax += 10
		case bool:
			tags[metricName] = strconv.FormatBool(metric)
			if metric {
				scoreValue += 1
			}
			scoreMax += 1
		default:
			return fmt.Errorf("unexpected metric type for %s, %v", metricName, metric)
		}
	}

	// aggregate tags and stringify values
	a.updates <- sources.Result{
		Time: req.Date,
		Tags: tags,
		Fields: map[string]interface{}{
			"notes":        req.Notes,
			"score_value":  scoreValue,
			"score_max":    scoreMax,
			"score_health": float64(scoreValue / scoreMax * 100),
		},
	}

	return nil
}

func (a API) Updates() <-chan sources.Result {
	return a.updates
}
