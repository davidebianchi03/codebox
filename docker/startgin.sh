#!/bin/bash
cd /codebox/bin

# export PATH=/root/.nvm/versions/node/v20.12.2/bin:$PATH

mkdir -p /codebox/db
mkdir -p /codebox/data

echo "Applying migrations..."
atlas migrate apply --url "${CODEBOX_DB_URL}"

echo "Starting codebox server..."

./codebox runserver &
