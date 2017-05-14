#!/bin/bash

# get global vars
. scripts/config.sh

# het image id (docker images)
LAST_IMAGE=$(docker images tfeel-trader:latest -q)

# tag it
docker tag $LAST_IMAGE "gcr.io/${GCLOUD_PROJECT}/${APP_NAME}"

# push it
gcloud docker -- push "gcr.io/${GCLOUD_PROJECT}/${APP_NAME}:latest"
