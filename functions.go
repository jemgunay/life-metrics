package functions

import (
	"net/http"

	"github.com/jemgunay/life-metrics/api"
	"github.com/jemgunay/life-metrics/config"
	"github.com/jemgunay/life-metrics/influx"
)

var apiHandler api.API

func init() {
	conf := config.New()
	influxRequester := influx.New(conf.Influx)
	apiHandler = api.New(influxRequester)
}

func DayLogHandler(w http.ResponseWriter, r *http.Request) {
	apiHandler.Handler(w, r)
}