#!/bin/bash#!/bin/bash#!/bin/bash

gcloud functions deploy DayLogHandler --runtime go116 --trigger-http --allow-unauthenticated
