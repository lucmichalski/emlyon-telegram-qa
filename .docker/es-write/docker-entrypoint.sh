#!/bin/sh

curl -s -XPUT -H "Content-Type: application/json" http://elastic:9200/_cluster/settings -d '{ "transient": { "cluster.routing.allocation.disk.threshold_enabled": false } }' | jq .
curl -s -XPUT -H "Content-Type: application/json" http://elastic:9200/_all/_settings -d '{"index.blocks.read_only_allow_delete": null}' | jq .
