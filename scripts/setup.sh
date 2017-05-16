#!/bin/bash

DIR="$(dirname "$0")"
. "${DIR}/config.sh"


gcloud beta spanner instances create ${GCLOUD_INSTANCE} \
  --config=regional-us-central1 \
  --description="${GCLOUD_INSTANCE} Instance" \
  --nodes=1

gcloud beta spanner databases create ${GCLOUD_DB} \
  --instance=${GCLOUD_INSTANCE}

# some gymnastics are required in order to parse a proper DDL in commandline
echo 'Loading DDL...'
echo 'NOTE: empty @type property warning on return protobuf message are OK'
DDL=`cat ${DIR}/store.ddl | tr -d '\n' | tr -d '\r' | tr -d '\t'`
IFS=';' read -ra LINES <<< "$DDL"
for SQL in "${LINES[@]}"; do
    # echo $SQL
    if [ ${#SQL} -ge 5 ]; then
      gcloud beta spanner databases ddl update ${GCLOUD_DB} --instance=${GCLOUD_INSTANCE} --ddl="$SQL"
    fi
done
