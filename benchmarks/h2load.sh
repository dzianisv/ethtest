#!/bin/bash
set -euo pipefail
# apt install -yq nghttp2-client
BASE64_USER_PASSWORD=$(echo "$NAME:$PASSWORD" | base64)
exec h2load -n 1000000 -c 100 -H "Content-Type: application/json" -H 'User-Agent: h2load' --data data.json  -H "Authorization: Basic $BASE64_USER_PASSWORD" "https://$HOST/"