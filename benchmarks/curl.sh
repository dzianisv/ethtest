#!/bin/sh

exec curl -H "Content-Type: application/json" -H 'User-Agent: curl-ethtest' -u "$NAME:$PASSWORD" -m 5 --data @data.json -X POST  "https://$HOST/"