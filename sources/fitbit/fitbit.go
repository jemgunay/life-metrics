package fitbit

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jemgunay/life-metrics/config"
	"github.com/jemgunay/life-metrics/sources"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

// Fitbit represents the Fitbit collection source.
type Fitbit struct {

}

// New initialises the Fitbit source and manages auth token refreshing.
func New(conf config.Config, exporter sources.Exporter) *Fitbit {
	m := &Fitbit{

	}

	return m
}

// Name returns the source name.
func (f *Fitbit) Name() string {
	return "fitbit"
}

// Collect enqueues a Fitbit collection request.
func (f *Fitbit) Collect(period sources.Period) {
	userID := "bagels"
	sleepURL := fmt.Sprintf("/1.2/user/%s/sleep/list.json", userID)
}

// State returns the Fitbit running state.
func (f *Fitbit) State() sources.StateSet {
	return sources.StateSet{}
}
