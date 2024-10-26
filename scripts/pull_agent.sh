#!/bin/bash

AGENT_LINK=https://github.com/davidebianchi03/codebox-agent/releases/download/2024.10.24.01/codebox-agent-linux-amd64
echo "Dowloading agent from ${AGENT_LINK}"
wget -O codebox-agent-linux-amd64 ${AGENT_LINK}