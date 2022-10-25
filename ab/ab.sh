#!/bin/sh

ab -n 1000000 -c 100 -T "application/json" -H 'User-Agent: ab-ethtest' -s 5 -p data.json -m POST  "https://$USER:$PASSWORD@$HOST/"