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
  --zone $GCLOUD_ZONE \
  --network "default" \
  --enable-cloud-logging \
  --enable-cloud-monitoring

# connect
gcloud container clusters get-credentials "${APP_NAME}-cluster" \
  --zone $GCLOUD_ZONE --project $GCLOUD_PROJECT

# configs
kubectl create configmap tsignal-config --from-file configmaps/us-west1.yaml


# populate twitter secrets
kubectl create secret generic tsignal-tw-key --from-literal=T_CONSUMER_KEY=$T_CONSUMER_KEY
kubectl create secret generic tsignal-tw-secret --from-literal=T_CONSUMER_SECRET=$T_CONSUMER_SECRET
kubectl create secret generic tsignal-tw-token --from-literal=T_ACCESS_TOKEN=$T_ACCESS_TOKEN
kubectl create secret generic tsignal-tw-access --from-literal=T_ACCESS_SECRET=$T_ACCESS_SECRET

# populate app secrets
kubectl create secret generic tsignal-gcloud-project --from-literal=GCLOUD_PROJECT=$GCLOUD_PROJECT
kubectl create secret generic tsignal-spanner-instance --from-literal=GCLOUD_INSTANCE=$GCLOUD_INSTANCE
kubectl create secret generic tsignal-spanner-db --from-literal=GCLOUD_DB=$GCLOUD_DB
kubectl create secret generic tsignal-sa --from-file auth.json

# deploy
# kubectl create -f deployments/tsignal.yaml
