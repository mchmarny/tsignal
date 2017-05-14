#!/bin/bash

# get global vars
. scripts/config.sh

docker build -t "${APP_NAME}:latest" ./ \
  --build-arg T_CONSUMER_KEY=$T_CONSUMER_KEY \
  --build-arg T_CONSUMER_SECRET=$T_CONSUMER_SECRET \
  --build-arg T_ACCESS_TOKEN=$T_ACCESS_TOKEN \
  --build-arg T_ACCESS_SECRET=$T_ACCESS_SECRET \
  --build-arg GCLOUD_PROJECT=$GCLOUD_PROJECT \
  --build-arg GCLOUD_INSTANCE=$GCLOUD_INSTANCE \
  --build-arg GCLOUD_DB=$GCLOUD_DB \
  --build-arg GOOGLE_APPLICATION_CREDENTIALS=$GOOGLE_APPLICATION_CREDENTIALS

# run
# docker run -i -t $(docker images tfeel-trader:latest -q)
