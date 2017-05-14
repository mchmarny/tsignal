#!/bin/bash

# get global vars
. scripts/config.sh

# create cluster
gcloud container clusters create "${APP_NAME}-cluster" \
  --project ${GCLOUD_PROJECT} \
  --machine-type "n1-standard-1" \
  --image-type "COS" \
  --disk-size "100" \
  --scopes default,cloud-platform,logging-write,monitoring-write \
  --num-nodes "1" \
  --network "default" \
  --enable-cloud-logging \
  --enable-cloud-monitoring
