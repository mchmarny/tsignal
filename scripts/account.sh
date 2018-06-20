#!/bin/bash

DIR="$(dirname "$0")"
. "${DIR}/config.sh"


echo "Checking if Service Account alredy created..."
SA=$(gcloud iam service-accounts list --format='value(EMAIL)' --filter="EMAIL:${GCLOUD_SA_EMAIL}")
if [ -z "${SA}" ]; then
  echo "Service Account not set, creating..."
  gcloud beta iam service-accounts create $GCLOUD_SA_NAME \
    --display-name="${APP_NAME} service account"

  echo "Creating service account key..."
  gcloud iam service-accounts keys create --iam-account $GCLOUD_SA_EMAIL \
    $GOOGLE_APPLICATION_CREDENTIALS
fi

echo "Creating service account bindings..."
gcloud projects add-iam-policy-binding $GCLOUD_PROJECT \
    --member "serviceAccount:${GCLOUD_SA_EMAIL}" \
    --role='roles/spanner.databaseUser'
