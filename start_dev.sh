#!/bin/bash

docker build -t emergencyupgrader-dev . -f Dockerfile.dev
docker run  --name=emergencyupgrader-dev -p 8081:8081 --rm -it -v $(pwd):/app emergencyupgrader-dev
