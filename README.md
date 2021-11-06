# Life Metrics

[![CircleCI](https://circleci.com/gh/jemgunay/life-metrics/tree/master.svg?style=svg)](https://circleci.com/gh/jemgunay/life-metrics/tree/master)

Life metrics is a service for collecting daily health metrics and scraping personal data from third party APIs. Collected data is persisted to InfluxDB to be visualised in Grafana.  

## Data Sources

* Day log
* Monzo (Eating Out)

## TODO

* Google scheduler -> cron -> refresh handlers
* Firebase for persisting oauth tokens on restart
* Need general API/web app authentication
* Sources
  * FitBit
  * Canlendar
  * GitHub commits
  * Calendar (events, birthdays)
  * Monzo
  * Spotify
  * Uber
  * Phone usage 


