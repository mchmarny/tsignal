#!/bin/bash

DIR="$(dirname "$0")"
. "${DIR}/config.sh"

gcloud beta spanner instances delete ${GCLOUD_INSTANCE}
