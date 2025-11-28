#!/usr/bin/env bash
set -euo pipefail

RUNNER_REPO_ID="${RUNNER_REPO_ID:-69007830}"
CI_API_V4_URL="${CI_API_V4_URL:-https://gitlab.com/api/v4}"

TOKEN=""
HEADER_NAME=""

if [[ -n "${CI_JOB_TOKEN:-}" ]]; then
  TOKEN="$CI_JOB_TOKEN"
  HEADER_NAME="JOB-TOKEN"
  echo "Using CI_JOB_TOKEN for authentication."
elif [[ -n "${PERSONAL_ACCESS_TOKEN:-}" ]]; then
  TOKEN="$PERSONAL_ACCESS_TOKEN"
  HEADER_NAME="PRIVATE-TOKEN"
  echo "Using PERSONAL_ACCESS_TOKEN for authentication."
else
  echo "Error: No authentication token found."
  echo "Set one of these environment variables before running:"
  echo "  export CI_JOB_TOKEN=<your-ci-token>"
  echo "  or"
  echo "  export PERSONAL_ACCESS_TOKEN=<your-personal-token>"
  exit 1
fi

if ! command -v jq >/dev/null || ! command -v curl >/dev/null; then
  echo "Installing missing dependencies..."
  sudo apt update -qq && sudo apt install -y jq curl
fi

latest_version=$(curl --header "${HEADER_NAME}: ${TOKEN}" \
  --silent \
  --location \
  "${CI_API_V4_URL}/projects/${RUNNER_REPO_ID}/repository/tags" \
  | jq -r '[.[]][0].name')

if [[ -z "$latest_version" || "$latest_version" == "null" ]]; then
  echo "Failed to get latest runner version"
  return 1
fi

echo "Latest version: ${latest_version}"

echo "${latest_version}" > recommended_runner_version.txt

echo "Recommended runner version saved to recommended_runner_version.txt"
