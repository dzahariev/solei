#!/usr/bin/env bash
set -e

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

POSTGRES_CONTAINER="pg-db-solei-verify"
POSTGRES_HOST="localhost"
POSTGRES_PORT="5432" 
POSTGRES_DB="postgres"
POSTGRES_USER="postgres"
POSTGRES_PASSWORD="postgres"
POSTGRES_SSL="disable"
POSTGRES_VERSION="17-alpine"

CONNECTION_STRING="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL}"

function clean() {
    echo "Remove DB container"
    docker rm --force ${POSTGRES_CONTAINER}
}

trap clean EXIT

set +e
echo "Create DB container ${POSTGRES_CONTAINER} ..."
docker run -d --name ${POSTGRES_CONTAINER} \
            -e POSTGRES_HOST=${POSTGRES_HOST} \
            -e POSTGRES_PORT=${POSTGRES_PORT} \
            -e POSTGRES_DB=${POSTGRES_DB} \
            -e POSTGRES_USER=${POSTGRES_USER} \
            -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
            -p ${POSTGRES_PORT}:${POSTGRES_PORT} \
            postgres:${POSTGRES_VERSION}

echo "Waiting for DB to start ..."
for i in {1..10}
do
    docker exec ${POSTGRES_CONTAINER} pg_isready -U "${POSTGRES_USER}" -h "${POSTGRES_HOST}" -p "${POSTGRES_PORT}" -d "${POSTGRES_DB}"
    if [ $? -eq 0 ]
    then
        DB_STARTED=true
        break
    fi
    sleep 1
done

if [ "${DB_STARTED}" != true ] ; then
    echo "DB in container ${POSTGRES_CONTAINER} was not started. Exiting."
    exit 1
fi
set -e

echo "Execute DB Migrations UP ..."
migrate -verbose -path ${ROOT_PATH}/migrations -database "${CONNECTION_STRING}" up

echo "Load initial content ..."
ls -d ${ROOT_PATH}/init/* | sort | xargs -I {} cat {} |\
    docker exec -i ${POSTGRES_CONTAINER} psql -U "${POSTGRES_USER}" -h "${POSTGRES_HOST}" -p "${POSTGRES_PORT}" -d "${POSTGRES_DB}"

echo "Migration version: $(migrate -path ${ROOT_PATH}/migrations -database "${CONNECTION_STRING}" version 2>&1)"

echo "Execute DB Migrations DOWN ..."
migrate -verbose -path ${ROOT_PATH}/migrations -database "${CONNECTION_STRING}" down 1

echo "Migration version: $(migrate -path ${ROOT_PATH}/migrations -database "${CONNECTION_STRING}" version 2>&1)"
