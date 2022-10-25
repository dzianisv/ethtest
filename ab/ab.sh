#!/bin/sh

exec ab -n 1000000 -c 100 -T "application/json" -H 'User-Agent: ab-ethtest' -A "$USER:$PASSWORD" -s 5 -p data.json -m POST  "https://$HOST/"