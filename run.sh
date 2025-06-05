#!/usr/bin/env bash
set -e

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

PRESERVE_DB=false
PRESERVE_KEYCLOAK=false
REUSE_DB=false
WITH_UI=true
REBUILD_UI=true

POSITIONAL=()
while [[ $# -gt 0 ]]
do

    key="$1"

    case ${key} in
        --preserve-db)
            PRESERVE_DB=true
            shift
        ;;
        --reuse-db)
            REUSE_DB=true
            shift
        ;;
        --preserve-keycloak)
            PRESERVE_KEYCLOAK=true
            shift
        ;;
        --reuse-keycloak)
            REUSE_KEYCLOAK=true
            shift
        ;;
        --no-ui)
            WITH_UI=false
            shift
        ;;
        --no-rebuild-ui)
            REBUILD_UI=false
            shift
        ;;
        --debug)
            DEBUG=true
            DEBUG_PORT=40000
            shift
        ;;
        --*)
            echo "Unknown flag ${1}"
            exit 1
        ;;
    esac
done
set -- "${POSITIONAL[@]}" 

set -a
. ./.env
set +a

POSTGRES_CONTAINER="solei-db"
POSTGRES_VERSION="17-alpine"

KEYCLOAK_CONTAINER="solei-keycloak"
KEYCLOAK_VERSION="26.2.3"
KEYCLOAK_PORT="8086"

CONNECTION_STRING="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

function clean() {

    if [[ ${DEBUG} == true ]]; then
       echo "Remove binary"
       rm  ${ROOT_PATH}/solei
    fi

    if [[ ${PRESERVE_DB} = false ]]; then
        echo "Remove DB container"
        docker rm --force ${POSTGRES_CONTAINER}
    else
        echo "Keeping DB container running"
    fi

    if [[ ${PRESERVE_KEYCLOAK} = false ]]; then
        echo "Remove Keycloak container"
        docker rm --force ${KEYCLOAK_CONTAINER}
    else
        echo "Keeping Keycloak container running"
    fi
}

trap clean EXIT

# Ensure public folder exists
mkdir ${ROOT_PATH}/public || true

if [[ ${REUSE_DB} = true ]]; then
    echo "DB is reused."
else
    set +e
    echo "Create DB container ${POSTGRES_CONTAINER} ..."
    docker run -d --name ${POSTGRES_CONTAINER} \
                -e POSTGRES_HOST=${DB_HOST} \
                -e POSTGRES_PORT=${DB_PORT} \
                -e POSTGRES_DB=${DB_NAME} \
                -e POSTGRES_USER=${DB_USER} \
                -e POSTGRES_PASSWORD=${DB_PASSWORD} \
                -p ${DB_PORT}:${DB_PORT} \
                postgres:${POSTGRES_VERSION}

    if [[ $? -ne 0 ]] ; then
        PRESERVE_DB=true
        exit 1
    fi

    echo "Waiting for DB to start ..."
    for i in {1..10}
    do
        docker exec ${POSTGRES_CONTAINER} pg_isready -U "${DB_USER}" -h "${DB_HOST}" -p "${DB_PORT}" -d "${DB_NAME}"
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

    echo "Execute DB Migrations ..."
    migrate -path ${ROOT_PATH}/db/migrations -database "${CONNECTION_STRING}" up

    echo "Load initial content ..."
    ls -d ${ROOT_PATH}/db/init/* | sort | xargs -I {} cat {} |\
        docker exec -i ${POSTGRES_CONTAINER} psql -U "${DB_USER}" -h "${DB_HOST}" -p "${DB_PORT}" -d "${DB_NAME}"
fi

echo "Migration version: $(migrate -path ${ROOT_PATH}/db/migrations -database "${CONNECTION_STRING}" version 2>&1)"

if [[ ${REUSE_KEYCLOAK} = true ]]; then
    echo "Keycloak is reused."
else
    echo "Create Keycloak container ${KEYCLOAK_CONTAINER} on ${KEYCLOAK_PORT} ..."
    docker run -d --name ${KEYCLOAK_CONTAINER} \
                -e KC_BOOTSTRAP_ADMIN_USERNAME=${KEYCLOAK_ADMIN} \
                -e KC_BOOTSTRAP_ADMIN_PASSWORD=${KEYCLOAK_ADMIN_PASSWORD} \
                -p ${KEYCLOAK_PORT}:${KEYCLOAK_PORT} \
                -v ./keycloak:/opt/keycloak/data/import \
                keycloak/keycloak:${KEYCLOAK_VERSION} \
                start-dev --http-port ${KEYCLOAK_PORT} --import-realm
fi

if [[ ${WITH_UI} = true ]]; then
    if [[ ${REBUILD_UI} = true ]]; then
        echo "Rebuilding UI ..."
        if [ ! -d "${ROOT_PATH}/ui/keycloak-js" ]; then
            echo "${ROOT_PATH}/ui/keycloak-js is missing. Fetching it ..."
            sh update-keycloak-js.sh
        else
            echo "${ROOT_PATH}/ui/keycloak-js is available, using the local copy."
        fi 
        rm -fR ${ROOT_PATH}/public/*
        cd ${ROOT_PATH}/ui
        rm -fR node_modules
        rm -f package-lock.json
        npm install
        npm run build
        mv ${ROOT_PATH}/ui/dist/* ${ROOT_PATH}/public/.
        cp ${ROOT_PATH}/ui/sap-ui-version.json ${ROOT_PATH}/public/resources/.
        mkdir -p ${ROOT_PATH}/public/resources/libs
        cp -r ${ROOT_PATH}/ui/keycloak-js ${ROOT_PATH}/public/resources/libs/.
        cd ${ROOT_PATH}
        echo "UI was rebuilded!"
    else
        echo "UI not rebuilded!"
    fi 
else
    echo "Starting without UI."
fi 

echo "Starting application ..."

if [[ ${DEBUG} == true ]]; then
    echo "Debug enabled on port ${DEBUG_PORT}"
    CGO_ENABLED=0 go build -gcflags="all=-N -l" .
    dlv --listen=:${DEBUG_PORT} --headless=true --api-version=2 exec solei
else
    go run main.go
fi
