package monzo

import (
	"time"

	"github.com/jemgunay/life-metrics/sources"
)

// Monzo TODO
type Monzo struct {

}

// New TODO
func New() *Monzo {
	return &Monzo{}
}

// Name TODO
func (m *Monzo) Name() string {
	return "monzo"
}

// Collect TODO
func (m *Monzo) Collect(start, end time.Time) ([]sources.Result, error) {
	return nil, nil
}

// Shutdown TODO
func (m *Monzo) Shutdown() {

}
