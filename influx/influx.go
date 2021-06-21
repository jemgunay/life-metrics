package influx

import (
	"log"

	"github.com/influxdata/influxdb-client-go/v2"
	influxdbapi "github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/jemgunay/life-metrics/sources"
)

type Requester struct {
	writeClient influxdbapi.WriteAPI
	readClient  influxdbapi.QueryAPI
}

func New(influxHost, authToken string) Requester {
	const (
		org       = "***REMOVED***"
		bucket    = "life-metrics"
	)
	client := influxdb2.NewClient(influxHost, authToken)
	return Requester{
		writeClient: client.WriteAPI(org, bucket),
		readClient:  client.QueryAPI(org),
	}
}

func (r Requester) Write(measurement string, results ...sources.Result) {
	log.Printf("%+v", results)
	// Create point using full params constructor
	for _, result := range results {
		p := influxdb2.NewPoint(
			measurement,
			result.Tags,
			result.Fields,
			result.Time,
		)

		r.writeClient.WritePoint(p)
	}
}
