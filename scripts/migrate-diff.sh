#!/bin/bash

set -e

script_dir=$(cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)
working_dir="$(dirname $script_dir)"
dotenv_file="${working_dir}/codebox.env"

echo "Loading config from $dotenv_file"
export $(grep -v '^#' ${dotenv_file} | xargs)

cd $working_dir && go get ariga.io/atlas-go-sdk/atlasexec && atlas migrate diff --env codebox
