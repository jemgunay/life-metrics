package daylog

import (
	"fmt"
	"time"

	"github.com/jemgunay/life-metrics/influx"
)

// Request represents a day log creation request body.
type Request struct {
	Date    time.Time `json:"date"`
	Metrics metrics   `json:"metrics"`
	Notes   string    `json:"notes"`
}

// Response represents a day log data Response.
type Response struct {
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

// DayLog defines the DayLog handler entry point and access to influx.
type DayLog struct {
	influxRequester influx.Requester
}

// New returns an initialised DayLog.
func New(influxRequester influx.Requester) DayLog {
	return DayLog{
		influxRequester: influxRequester,
	}
}

// Fetch reads the day log for the given date.
func (d DayLog) Fetch(date time.Time) (Response, error) {
	data, err := d.influxRequester.ReadDayLog(date)
	if err != nil {
		return Response{}, fmt.Errorf("failed to query influx: %s", err)
	}

	// process influx data into response expected by the UI
	dayLogResp := Response{}
	if len(data) > 0 {
		dayLogResp.Submitted = true
		if data["notes"] != nil {
			dayLogResp.Notes = data["notes"].(string)
		}
		// remove notes as we have it covered above
		delete(data, "notes")
		dayLogResp.Metrics = data
	}

	return dayLogResp, nil
}

// Submit is the root HTTP DayLog handler for submitting and reading day logs.
func (d DayLog) Submit(req Request) error {
	if err := d.processDayLog(req); err != nil {
		return fmt.Errorf("failed to process Request data: %s", err)
	}
	return nil
}

// processDayLog processes the day log Request into a result to be written to influx.
func (d DayLog) processDayLog(req Request) error {
	res := newFieldSet(req.Notes)

	// add all Request fields to the result
	res.addInt("general_mood", req.Metrics.GeneralMood)
	res.addInt("diet_quality", req.Metrics.DietQuality)
	res.addInt("water_intake", req.Metrics.WaterIntake)
	res.addIntInvert("caffeine_intake", req.Metrics.CaffeineIntake)
	res.addBool("exercise", req.Metrics.Exercise)
	res.addBool("meditation", req.Metrics.Meditation)
	res.calcHealth()

	// aggregate tags and stringify values
	logData := influx.Result{
		Time:   req.Date,
		Fields: res.fields,
	}

	if err := d.influxRequester.Write("day_log", logData); err != nil {
		return fmt.Errorf("failed to write day log data to influx: %s", err)
	}
	return nil
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

// addInt adds an int field to the fieldSet, updating the score (maximum of 10).
func (f *fieldSet) addInt(name string, value int) {
	f.fields[name] = value
	f.scoreValue += value
	f.scoreMax += 10
}

// addIntInvert subtracts the int from the max possible and adds the int field to the fieldSet, updating the score
// (maximum of 10).
func (f *fieldSet) addIntInvert(name string, value int) {
	f.fields[name] = 10 - value
	f.scoreValue += 10 - value
	f.scoreMax += 10
}

// addBool adds a bool field to the fieldSet, updating the score (maps false/true to 0/5).
func (f *fieldSet) addBool(name string, value bool) {
	f.fields[name] = false
	if value {
		f.fields[name] = true
		f.scoreValue += 5
	}
	f.scoreMax += 5
}

// calcHealth calculates the health score from the total and max field scores and adds it to the fieldSet.
func (f *fieldSet) calcHealth() {
	f.fields["score_value"] = f.scoreValue
	f.fields["score_max"] = f.scoreMax
	f.fields["score_health"] = float64(f.scoreValue) / float64(f.scoreMax) * 100
}
