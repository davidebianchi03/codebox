#!/bin/bash

REPO_OWNER=davidebianchi03
REPO_NAME=codebox-agent
LINUX_AGENT_ARTIFACT_NAME=codebox-agent-linux-amd64

response=$(curl -s -H "Authorization: token $GH_TOKEN" "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases")
releases_created_at=$(echo "$response" | jq '.[].created_at')
sorted_releases=$(echo $releases_created_at | xargs -n1 | sort -r | xargs)
latest_release_created_at=$(echo $sorted_releases | cut -d ' ' -f1)

latest_release=$(echo "$response" | jq ".[] | select(.created_at == \"${latest_release_created_at}\")")
linux_agent_asset_id=$(echo $latest_release | jq ".assets[] | select(.name == \"${LINUX_AGENT_ARTIFACT_NAME}\") | .id")

curl -L -H "Authorization: token $GH_TOKEN" -H "Accept: application/octet-stream" "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/assets/$linux_agent_asset_id" -o agent
