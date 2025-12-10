#!/bin/bash

# Ensure a migration name is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <migration_name>"
  exit 1
fi

MIGRATION_NAME="$1"

working_dir="$(dirname $(cd "$(dirname "$0")" && pwd))"
dotenv_file="${working_dir}/codebox.env"

export $(grep -v '^#' ${dotenv_file} | xargs)

cd $working_dir && \
  go get ariga.io/atlas-go-sdk/atlasexec && \
  atlas migrate new --env codebox "$MIGRATION_NAME"
  