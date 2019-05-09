#!/usr/bin/env bash

rm -f fairway

GOOS=linux GOARCH=amd64 go build -o fairway
docker build . -t fairway-docker --no-cache