package sources

import "time"

// Source defines the requirements for a data collection source.
type Source interface {
	// Name returns the name of the source.
	Name() string
	// Collect collects data for a given start and end time.
	Collect(start, end time.Time) ([]Result, error)
	// Shutdown shuts down the source.
	Shutdown()
}

// Result is a generic collection dataset returned from a source.
type Result struct {
	Time   time.Time
	Tags   map[string]string
	Fields map[string]interface{}
}
