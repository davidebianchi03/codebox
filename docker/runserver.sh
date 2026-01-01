#!/bin/bash

echo "Setting up environment..."

if [ -z "$CODEBOX_DB_HOST" ]; then
    echo "CODEBOX_DB_HOST not set, using default value";
    export CODEBOX_DB_HOST=db;
fi

if [ -z "$CODEBOX_DB_PORT" ]; then
    echo "CODEBOX_DB_PORT not set, using default value";
    export CODEBOX_DB_HOST=3306;
fi

if [ -z "$CODEBOX_DB_NAME" ]; then
    echo "CODEBOX_DB_NAME not set, using default value";
    export CODEBOX_DB_HOST=codebox;
fi

if [ -z "$CODEBOX_DB_USER" ]; then
    echo "CODEBOX_DB_USER not set, using default value";
    export CODEBOX_DB_HOST=codebox;
fi

if [ -z "$CODEBOX_DB_PASSWORD" ]; then
    echo "CODEBOX_DB_PASSWORD not set, using default value";
    export CODEBOX_DB_HOST=password;
fi

cd /codebox/bin

mkdir -p /codebox/data

export $(grep -v '^#' codebox.env | xargs)

echo "Applying migrations..."
atlas migrate apply --env codebox

echo "Starting nginx..."
nginx -g "daemon off;" > /dev/null 2>&1 &

echo "Starting codebox..."
./codebox runserver
