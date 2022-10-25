#!/bin/sh

exec h2load -n 1000000 -c 100 -H "Content-Type: application/json" -H 'User-Agent: h2load' --data data.json  "https://$NAME:$PASSWORD@$HOST/"