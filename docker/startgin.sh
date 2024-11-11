#!/bin/bash
cd /codebox/bin

mkdir -p /codebox/db
mkdir -p /codebox/data

./codebox runserver &
