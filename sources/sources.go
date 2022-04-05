package sources

import "time"

// StateSet represents the current running state of a source - this is used by the sources page.
type StateSet map[string]interface{}

// Exporter defines the contract for writing data to be persisted outside of the service, e.g. InfluxDB.
type Exporter interface {
	Write(measurement string, results ...Result) error
}

// Result is a generic collection dataset returned from a source.
type Result struct {
	Time   time.Time
	Tags   map[string]string
	Fields map[string]interface{}
}

// Period represents a time period.
type Period struct {
	Start time.Time
	End   time.Time
}

// NewPeriod creates a new Period given a start and end time.
func NewPeriod(start, end time.Time) Period {
	return Period{
		Start: start,
		End:   end,
	}
}
