#!/bin/bash
cd /codebox/bin

mkdir -p /codebox/db
mkdir -p /codebox/data

export $(grep -v '^#' codebox.env | xargs)

echo "Applying migrations..."
atlas migrate apply --env codebox

./codebox runserver &
