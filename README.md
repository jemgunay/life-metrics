# Life Metrics

[![CircleCI](https://circleci.com/gh/jemgunay/life-metrics/tree/master.svg?style=svg)](https://circleci.com/gh/jemgunay/life-metrics/tree/master)

Life metrics is a service for collecting daily health metrics and scraping personal data from third party APIs. Collected data is persisted to InfluxDB and visualised in Grafana.  

<p align="center">
  <img src="/screenshots/screenshot_1.png" width="30%"/>
  <img src="/screenshots/screenshot_2.png" width="30%"/>
  <img src="/screenshots/screenshot_3.png" width="30%"/>
</p>

## Run

Configure the env vars defined in `config/env-setup.sh` then run:

```bash
source config/env-setup.sh
make local
go run life-metrics.go
```

## Data Sources

### Implemented Sources

* Day log form
* Monzo ("Eating Out" category transactions)

### Day Log Endpoint

The day log endpoint is responsible for submitting day log data and for retrieving submitted day logs.

* Fetch day log data for a given date:
```bash
curl -i "http//localhost:8080/api/data/daylog?date=2021-11-07T00:00:00Z" -XGET
```

* Submit a date's day log:
```bash
curl -i "http//localhost:8080/api/data/daylog" -XPOST -d '{"date":"2021-11-07T00:00:00Z","notes":"","metrics":{"general_mood":7,"diet_quality":3,"water_intake":4,"caffeine_intake":0,"exercise":false,"meditation":false}}'
```

### Collect Endpoint

The collect endpoint triggers a data collection for all sources. Collected data is then written to InfluxDB. Collection is triggered every 6 hours by the Google Cloud Scheduler.  

* Collection between the current time and the timestamp for the last series written by each source   
```bash
curl -i "http://localhost:8080/api/data/collect" -XPOST
```

* Collection for all data that each source can provide (overwrites existing data in Influx)
```bash
curl -i "http://localhost:8080/api/data/collect?reset=true" -XPOST
```

### Auth Endpoints

OAuth2 authentication endpoints:

* `/api/auth/monzo`

## TODO

* Add Vue CI lint & build
* Refresh per source API & ability to refresh for specified time range
* Firebase for persisting OAuth tokens on restart
* General API/web app authentication
* Sources
  * FitBit - exercise sessions & sleep data 
  * Canlendar - get alcohol units consumed into life-metrics and re-point canlendar service
  * Spotify - correlate genres/playlists with mood
  * Phone screen time - app required