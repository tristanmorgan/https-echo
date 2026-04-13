#!/bin/sh

export DOCKER_HOST=ssh://tristan@inspirone.node.home.consul

TS_VAR=$(date +%s)
docker build --pull --progress plain --no-cache --platform=linux/arm64/v8 -t tristanmorgan/https-echo:$TS_VAR-arm64 .
docker push tristanmorgan/https-echo:$TS_VAR-arm64
docker build --pull --progress plain --no-cache --platform=linux/amd64 -t tristanmorgan/https-echo:$TS_VAR-amd64 .
docker push tristanmorgan/https-echo:$TS_VAR-amd64

docker manifest create tristanmorgan/https-echo:$TS_VAR --amend tristanmorgan/https-echo:$TS_VAR-arm64 --amend tristanmorgan/https-echo:$TS_VAR-amd64
docker manifest push tristanmorgan/https-echo:$TS_VAR
docker manifest rm tristanmorgan/https-echo:$TS_VAR

docker manifest create tristanmorgan/https-echo:latest --amend tristanmorgan/https-echo:$TS_VAR-arm64 --amend tristanmorgan/https-echo:$TS_VAR-amd64
docker manifest push tristanmorgan/https-echo:latest
docker manifest rm tristanmorgan/https-echo:latest
