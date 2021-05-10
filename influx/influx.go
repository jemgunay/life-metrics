package influx

import (
	"log"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	influxdbapi "github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/jemgunay/life-metrics/sources"
)

type Requester struct {
	writeClient influxdbapi.WriteAPI
	readClient  influxdbapi.QueryAPI
}

func New() Requester {
	const (
		org       = "jem"
		bucket    = "life-metrics"
		authToken = ""
	)
	client := influxdb2.NewClient("http://localhost:8086", authToken)
	return Requester{
		writeClient: client.WriteAPI(org, bucket),
		readClient:  client.QueryAPI(org),
	}
}

func (r Requester) Write(results []sources.Result) {
	// TODO: write to influx
	log.Printf("%+v", results)
	// Create point using full params constructor
	p := influxdb2.NewPoint(
		"stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45.0},
		time.Now().UTC(),
	)

	r.writeClient.WritePoint(p)
}
