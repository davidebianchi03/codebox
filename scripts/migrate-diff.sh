#!/bin/bash
working_dir="$(dirname $(cd "$(dirname "$0")" && pwd))"
dotenv_file="${working_dir}/codebox.env"
export $(grep -v '^#' ${dotenv_file} | xargs)

cd $working_dir && atlas migrate diff --env codebox
