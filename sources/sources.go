package sources

import "time"

type Source interface {
	Name() string
	Collect(start, end time.Time) ([]Result, error)
	Shutdown()
}

type Result struct {
	Time time.Time
	Value float64
	Metadata interface{}
}
