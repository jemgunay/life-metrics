package influx

import (
	"context"
	"fmt"
	"log"

	"github.com/influxdata/influxdb-client-go/v2"
	influxdbapi "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/jemgunay/life-metrics/sources"
)

type Requester struct {
	writeClient influxdbapi.WriteAPIBlocking
	readClient  influxdbapi.QueryAPI
}

func New(influxHost, authToken string) Requester {
	const (
		org    = "jemgunay@gmail.com"
		bucket = "life-metrics"
	)
	client := influxdb2.NewClient(influxHost, authToken)
	return Requester{
		writeClient: client.WriteAPIBlocking(org, bucket),
		readClient:  client.QueryAPI(org),
	}
}

func (r Requester) Write(measurement string, results ...sources.Result) error {
	log.Printf("writing to influx: %+v", results)

	points := make([]*write.Point, 0, len(results))
	for _, result := range results {
		point := influxdb2.NewPoint(
			measurement,
			result.Tags,
			result.Fields,
			result.Time,
		)
		points = append(points, point)
	}

	if err := r.writeClient.WritePoint(context.Background(), points...); err != nil {
		return fmt.Errorf("writing points to influx failed: %s", err)
	}

	return nil
}
