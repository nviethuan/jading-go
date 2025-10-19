#!/bin/bash

rm -rf ./fsell

docker buildx build -f ./Dockerfile.fsell -t my-app-fsell .
docker create --name my-container my-app-fsell
docker cp my-container:/out/fsell ./
docker rm my-container
docker rmi my-app-fsell

echo "Copying main file to $JADINGIP"
scp ~/Sources/mine-project/jading-go/fsell ubuntu@$JADINGIP:/home/ubuntu/app
echo "Copying completed"
