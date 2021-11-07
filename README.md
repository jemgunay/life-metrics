# Life Metrics

[![CircleCI](https://circleci.com/gh/jemgunay/life-metrics/tree/master.svg?style=svg)](https://circleci.com/gh/jemgunay/life-metrics/tree/master)

Life metrics is a service for collecting daily health metrics and scraping personal data from third party APIs. Collected data is persisted to InfluxDB and visualised in Grafana.  

## Data Sources

### Collect Endpoint

The collect endpoint triggers a data collection for all sources. Collected data is then written to InfluxDB. Collection is triggered every 6 hours by the Google Cloud Scheduler.  

* Collection between the current time and the timestamp for the last series written by each source   
```bash
curl -i http://localhost:8080/api/data/collect
```

* Collection for all data that each source can provide (overwrites existing data in Influx)
```bash
curl -i http://localhost:8080/api/data/collect?reset=true
```

### Implemented Sources

* Day log form
* Monzo ("Eating Out" category)

## TODO

* Refresh per source API & ability to refresh for specified time range
* Firebase for persisting OAuth tokens on restart
* Need general API/web app authentication
* Source refresh btn
* Sources
  * FitBit - exercise sessions & sleep data 
  * Canlendar - get alcohol units consumed into life-metrics and re-point canlendar service
  * Spotify - correlate genres/playlists with mood
  * Phone screen time - app required