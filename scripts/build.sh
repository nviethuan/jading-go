#!/bin/bash

rm -rf main

docker buildx build -t my-app .
docker create --name my-container my-app
docker cp my-container:/out/main ./
docker rm my-container
docker rmi my-app

echo "Copying main file to $JADINGIP"
scp ~/Sources/mine-project/jading-go/main ubuntu@$JADINGIP:/home/ubuntu/app
echo "Copying completed"