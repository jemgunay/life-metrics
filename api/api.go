package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jemgunay/life-metrics/influx"
	"github.com/jemgunay/life-metrics/sources"
)

// dayLogRequest represents a day log creation request body.
type dayLogRequest struct {
	Date    time.Time `json:"date"`
	Metrics metrics   `json:"metrics"`
	Notes   string    `json:"notes"`
}

// dayLogResponse represents a day log data response.
type dayLogResponse struct {
	Submitted bool                   `json:"submitted"`
	Metrics   map[string]interface{} `json:"metrics,omitempty"`
	Notes     string                 `json:"notes,omitempty"`
}

type metrics struct {
	GeneralMood    int  `json:"general_mood"`
	DietQuality    int  `json:"diet_quality"`
	WaterIntake    int  `json:"water_intake"`
	CaffeineIntake int  `json:"caffeine_intake"`
	Exercise       bool `json:"exercise"`
	Meditation     bool `json:"meditation"`
}

// API defines the API handler entry point and access to influx.
type API struct {
	influxRequester influx.Requester
}

// New returns an initialised API.
func New(influxRequester influx.Requester) API {
	return API{
		influxRequester: influxRequester,
	}
}

// Handler is the root HTTP API handler for submitting and reading day logs.
func (a API) Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request to API handler [%s] (%s) from %s", r.Method, r.URL, r.RemoteAddr)

	switch r.Method {
	// get today's submitted day log data
	case http.MethodGet:
		data, err := a.influxRequester.ReadDayLog(time.Now())
		if err != nil {
			log.Printf("failed to query influx: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// process influx data into response expected by the UI
		dayLogResp := dayLogResponse{}
		if len(data) > 0 {
			dayLogResp.Submitted = true
			if data["notes"] != nil {
				dayLogResp.Notes = data["notes"].(string)
			}
			// remove notes as we have it covered above
			delete(data, "notes")
			dayLogResp.Metrics = data
		}

		body, err := json.Marshal(dayLogResp)
		if err != nil {
			log.Printf("failed to JSON encode response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(body)

	case http.MethodPost:
		// submit today's day log data
		logReq, err := decodeBody(r)
		if err != nil {
			log.Printf("failed to process request body: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := a.processDayLog(logReq); err != nil {
			log.Printf("failed to process request data: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func decodeBody(r *http.Request) (dayLogRequest, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return dayLogRequest{}, fmt.Errorf("failed to read body: %s", err)
	}

	req := dayLogRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		return dayLogRequest{}, fmt.Errorf("failed to JSON decode body: %s", err)
	}

	return req, nil
}

// fieldSet is used to store fields to be written to influx and is used to calculate the score for a day log.
type fieldSet struct {
	fields     map[string]interface{}
	scoreValue int
	scoreMax   int
}

// newFieldSet creates an initialised fieldSet.
func newFieldSet(notes string) fieldSet {
	return fieldSet{
		fields: map[string]interface{}{
			"notes":           notes,
			"submission_date": time.Now().UTC(),
		},
	}
}

// addInt adds an int field to the fieldSet, updating the score.
func (f *fieldSet) addInt(name string, value int) {
	f.fields[name] = value
	f.scoreValue += value
	f.scoreMax += 10
}

// addBool adds a bool field to the fieldSet, updating the score.
func (f *fieldSet) addBool(name string, value bool) {
	f.fields[name] = false
	if value {
		f.fields[name] = true
		f.scoreValue++
	}
	f.scoreMax++
}

// calcHealth calculates the health score from the total and max field scores and adds it to the fieldSet.
func (f *fieldSet) calcHealth() {
	f.fields["score_value"] = f.scoreValue
	f.fields["score_max"] = f.scoreMax
	f.fields["score_health"] = float64(f.scoreValue) / float64(f.scoreMax) * 100
}

// processDayLog processes the day log request into a result to be written to influx.
func (a API) processDayLog(req dayLogRequest) error {
	res := newFieldSet(req.Notes)

	// add all request fields to the result
	res.addInt("general_mood", req.Metrics.GeneralMood)
	res.addInt("diet_quality", req.Metrics.DietQuality)
	res.addInt("water_intake", req.Metrics.WaterIntake)
	res.addInt("caffeine_intake", req.Metrics.CaffeineIntake)
	res.addBool("exercise", req.Metrics.Exercise)
	res.addBool("meditation", req.Metrics.Meditation)
	res.calcHealth()

	// aggregate tags and stringify values
	logData := sources.Result{
		Time:   req.Date,
		Fields: res.fields,
	}

	if err := a.influxRequester.Write("day_log", logData); err != nil {
		return fmt.Errorf("failed to write day log data to influx: %s", err)
	}
	return nil
}
