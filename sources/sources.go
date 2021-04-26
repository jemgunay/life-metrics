package sources

import "time"

type Source interface {
	Collect(start, end time.Time) ([]Result, error)
	Shutdown()
}

type Result struct {
	Time time.Time
	Value float64
	Metadata interface{}
}


