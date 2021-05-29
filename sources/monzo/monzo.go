package monzo

import (
	"time"

	"github.com/jemgunay/life-metrics/config"
	"github.com/jemgunay/life-metrics/sources"
)

type Monzo struct {

}

func New(conf config.Monzo) *Monzo {
	return &Monzo{}
}

func (m *Monzo) Name() string {
	return "monzo"
}

func (m *Monzo) Collect(start, end time.Time) ([]sources.Result, error) {
	return nil, nil
}

func (m *Monzo) Shutdown() {

}
