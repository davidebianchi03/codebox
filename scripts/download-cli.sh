#!/usr/bin/env bash
set -euo pipefail

CLI_REPO_ID="${CLI_REPO_ID:-68940749}"
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

rm -rf bin/cli/*
mkdir -p bin/cli
cd bin/cli

download_latest_package() {
  local package_name="$1"
  local file_name="$2"
  local output_file_name="$3"

  echo "Fetching latest version for ${package_name}..."
  local latest_version
  latest_version=$(curl --header "${HEADER_NAME}: ${TOKEN}" \
    --silent \
    --location \
    "${CI_API_V4_URL}/projects/${CLI_REPO_ID}/packages?package_name=${package_name}&package_type=generic&order_by=created_at&sort=desc" \
    | jq -r '.[0].version')

  if [[ -z "$latest_version" || "$latest_version" == "null" ]]; then
    echo "Failed to get latest version for ${package_name}"
    return 1
  fi

  echo "Latest version: ${latest_version}"
  echo "Downloading ${file_name}..."

  curl --header "${HEADER_NAME}: ${TOKEN}" \
       --location \
       --output "${output_file_name}" \
       "${CI_API_V4_URL}/projects/${CLI_REPO_ID}/packages/generic/${package_name}/${latest_version}/${file_name}"

  chmod +x "${output_file_name}" || true
  echo "Saved to: $(pwd)/${output_file_name}"
  echo
}

declare -a ARCH=("amd64" "386" "arm" "arm64")
for a in "${ARCH[@]}"; do
    echo "Processing architecture $a"

    if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "arm64" ]; then
      declare -a FILES=("codebox-cli-linux-${a}" "codebox-cli-${a}.deb" "codebox-cli-windows-${a}.exe" "codebox-cli-setup-${a}.exe" "codebox-cli-darwin-${a}")
    else
      declare -a FILES=("codebox-cli-linux-${a}" "codebox-cli-${a}.deb" "codebox-cli-windows-${a}.exe" "codebox-cli-setup-${a}.exe")
    fi

    for f in "${FILES[@]}"; do
      echo "Downloading $f"
      download_latest_package codebox-cli "${f}" "${f}"
    done

done
