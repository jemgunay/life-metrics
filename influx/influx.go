package influx

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	influxdbapi "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/jemgunay/life-metrics/sources"
)

const (
	org    = "***REMOVED***"
	bucket = "life-metrics"
)

// Requester is used to write to and query influx.
type Requester struct {
	writeClient influxdbapi.WriteAPIBlocking
	readClient  influxdbapi.QueryAPI
}

// New returns an initialised influx requester.
func New(influxHost, authToken string) Requester {
	client := influxdb2.NewClient(influxHost, authToken)
	return Requester{
		writeClient: client.WriteAPIBlocking(org, bucket),
		readClient:  client.QueryAPI(org),
	}
}

// Write writes the provided data to influx.
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

// ReadDayLog queries influx for the current day log's metrics.
func (r Requester) ReadDayLog(day time.Time) (map[string]interface{}, error) {
	startTime := day.Truncate(time.Hour * 24)
	endTime := startTime.Add(time.Hour * 24).Add(-time.Second)

	query := `from(bucket: "life-metrics")
  	|> range(start: ` + startTime.Format(time.RFC3339) + `, stop: ` + endTime.Format(time.RFC3339) + `)
  	|> filter(fn:(r) =>
    	r._measurement == "day_log"
  	)
  	|> last()`

	result, err := r.readClient.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query influx: %s", err)
	}

	data := make(map[string]interface{})
	for result.Next() {
		data[result.Record().Field()] = result.Record().Value()
	}

	return data, nil
}
