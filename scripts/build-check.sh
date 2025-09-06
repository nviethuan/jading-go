#!/bin/bash

rm -rf check/main

docker buildx build -f ./Dockerfile.check -t my-app-check .
docker create --name my-container my-app-check
docker cp my-container:/out/main ./check
docker rm my-container
docker rmi my-app-check
