#!/bin/bash

read -rp "Enter influx host: " influx_host
read -rsp "Enter influx token: " influx_token

curl --request POST "${influx_host}/api/v2/delete/\?org\=jemgunay\@gmail\.com/&bucket/=life-metrics" \
  --header "Authorization: Token ${influx_token}" \
  --header 'Content-Type: application/json' \
  --data '{
    "start": "2020-03-01T00:00:00Z",
    "stop": "2021-11-14T00:00:00Z"
  }'

if [[ $? != 0 ]]; then
  printf "\nInflux purge failed\n"
  exit 1
fi

printf "\nInflux purged\n"