package influx

import (
	"log"

	"github.com/jemgunay/life-metrics/sources"
)

type Requester struct {

}

func New() Requester {
	return Requester{}
}

func (r Requester) Write(results []sources.Result) {
	// TODO: write to influx
	log.Printf("%+v", results)
}