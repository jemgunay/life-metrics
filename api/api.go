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

type request struct {
	Date    time.Time `json:"date"`
	Metrics metrics   `json:"metrics"`
	Notes   string    `json:"notes"`
}

type metrics struct {
	GeneralMood  int  `json:"general_mood"`
	WorkMood     int  `json:"work_mood"`
	DietQuality  int  `json:"diet_quality"`
	WaterIntake  int  `json:"water_intake"`
	SleepQuality int  `json:"sleep_quality"`
	Exercise     bool `json:"exercise"`
	Meditation   bool `json:"meditation"`
}

type API struct {
	influxRequester influx.Requester
}

func New(influxRequester influx.Requester) API {
	return API{
		influxRequester: influxRequester,
	}
}

func (a API) Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request to API handler (%s) from %s", r.URL, r.RemoteAddr)

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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("request to API handler processed")
}

type result struct {
	fields     map[string]interface{}
	scoreValue int
	scoreMax   int
}

func newResult(notes string) result {
	return result{
		fields: map[string]interface{}{
			"notes":           notes,
			"submission_date": time.Now(),
		},
	}
}

func (r *result) calcHealth() {
	r.fields["score_value"] = r.scoreValue
	r.fields["score_max"] = r.scoreMax
	r.fields["score_health"] = float64(r.scoreValue) / float64(r.scoreMax) * 100
}

func (r *result) addInt(name string, value int) {
	r.fields[name] = value
	r.scoreValue += value
	r.scoreMax += 10
}

func (r *result) addBool(name string, value bool) {
	r.fields[name] = 0
	if value {
		r.fields[name] = 1
		r.scoreValue += 1
	}
	r.scoreMax += 1
}

func (a API) process(req request) error {
	res := newResult(req.Notes)

	// add all request fields to the result
	res.addInt("general_mood", req.Metrics.GeneralMood)
	res.addInt("work_mood", req.Metrics.WorkMood)
	res.addInt("diet_quality", req.Metrics.DietQuality)
	res.addInt("water_intake", req.Metrics.WaterIntake)
	res.addInt("sleep_quality", req.Metrics.SleepQuality)
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
