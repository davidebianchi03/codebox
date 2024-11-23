#!/bin/bash
cd /codebox/bin

export PATH=/root/.nvm/versions/node/v20.12.2/bin:$PATH

mkdir -p /codebox/db
mkdir -p /codebox/data

./codebox runserver &
