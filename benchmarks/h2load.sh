#!/bin/bash
# apt install -yq nghttp2-client
exec h2load -n 1000000 -c 100 -H "Content-Type: application/json" -H 'User-Agent: h2load'  -H "Authorization: Basic $(echo -n $NAME:$PASSWORD | base64 | tr -d '\n')" --data data.json "https://$HOST/"