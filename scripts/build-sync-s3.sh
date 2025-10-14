#!/bin/bash

rm -rf tools/ss3

docker buildx build -f ./Dockerfile.sync-s3 -t my-app-sync-s3 .
docker create --name my-container my-app-sync-s3
docker cp my-container:/out/ss3 ./tools
docker rm my-container
docker rmi my-app-sync-s3
