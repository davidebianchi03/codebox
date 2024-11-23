#!/bin/bash

REPO_OWNER=davidebianchi03
REPO_NAME=codebox-cli
LINUX_CLI_ARTIFACT_NAME=codebox-cli-linux-amd64
WINDOWS_CLI_ARTIFACT_NAME=codebox-cli-windows-amd64.exe


response=$(curl -s -H "Authorization: token $GH_TOKEN" "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases")
releases_created_at=$(echo "$response" | jq '.[].created_at')
sorted_releases=$(echo $releases_created_at | xargs -n1 | sort -r | xargs)
latest_release_created_at=$(echo $sorted_releases | cut -d ' ' -f1)

latest_release=$(echo "$response" | jq ".[] | select(.created_at == \"${latest_release_created_at}\")")
echo "CLI Latest Release: ${latest_release}"

linux_cli_asset_id=$(echo $latest_release | jq ".assets[] | select(.name == \"${LINUX_CLI_ARTIFACT_NAME}\") | .id")
echo "CLI Linux asset id: ${linux_cli_asset_id}"

windows_cli_asset_id=$(echo $latest_release | jq ".assets[] | select(.name == \"${WINDOWS_CLI_ARTIFACT_NAME}\") | .id")
echo "CLI Windows asset id: ${windows_cli_asset_id}"

curl -L -H "Authorization: token $GH_TOKEN" -H "Accept: application/octet-stream" "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/assets/$linux_cli_asset_id" -o $LINUX_CLI_ARTIFACT_NAME
curl -L -H "Authorization: token $GH_TOKEN" -H "Accept: application/octet-stream" "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/assets/$windows_cli_asset_id" -o $WINDOWS_CLI_ARTIFACT_NAME
