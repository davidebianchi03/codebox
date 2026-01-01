#!/bin/bash
working_dir="$(dirname $(cd "$(dirname "$0")" && pwd))"
dotenv_file="${working_dir}/codebox.env"
export $(grep -v '^#' ${dotenv_file} | xargs)
CODEBOX_DB_NAME=${CODEBOX_TEST_DB_NAME}
cd $working_dir

# remove all tables
atlas schema clean --env codebox --auto-approve

# migrate test db
atlas migrate apply --env codebox
