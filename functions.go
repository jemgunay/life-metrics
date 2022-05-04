package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"

	"github.com/jemgunay/life-metrics/config"
	"github.com/jemgunay/life-metrics/daylog"
	"github.com/jemgunay/life-metrics/influx"
)

var (
	apiHandler daylog.DayLog
	Conf = config.New()
)

func init() {
	influxRequester := influx.New(Conf.Influx)
	apiHandler = daylog.New(influxRequester)
}

func DayLogHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request to DayLog handler [%s] (%s) from %s", r.Method, r.URL, r.RemoteAddr)

	enableCORS(w)

	switch r.Method {
	// get today's submitted day log data
	case http.MethodGet:
		date, err := extractDateQuery(r)
		if err != nil {
			log.Printf("failed to process date query: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dayLogResp, err := apiHandler.Fetch(date)
		if err != nil {
			log.Printf("failed to fetch day log: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(dayLogResp)
		if err != nil {
			log.Printf("failed to JSON encode day log resp: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(body)

	// submit today's day log data
	case http.MethodPost:
		logReq, err := decodeBody(r)
		if err != nil {
			log.Printf("failed to process request body: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := apiHandler.Submit(logReq); err != nil {
			log.Printf("failed to fetch day log: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// enableCORS enables CORS for handlers that it wraps.
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func extractDateQuery(r *http.Request) (time.Time, error) {
	date := r.URL.Query().Get("date")
	if date == "" {
		return time.Time{}, errors.New("no date query provided")
	}

	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date as RFC3339: %s", err)
	}

	return parsedDate, nil
}

func decodeBody(r *http.Request) (daylog.Request, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return daylog.Request{}, fmt.Errorf("failed to read body: %s", err)
	}

	req := daylog.Request{}
	if err := json.Unmarshal(body, &req); err != nil {
		return daylog.Request{}, fmt.Errorf("failed to JSON decode body: %s", err)
	}

	if req.Date.IsZero() {
		return daylog.Request{}, fmt.Errorf("invalid date provided: %s", req.Date)
	}

	return req, nil
}
