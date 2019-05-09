#!/usr/bin/env bash

docker rm   fairway-docker
docker run -p 8001:8001 --name fairway-docker fairway-docker