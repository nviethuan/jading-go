#!/bin/bash

echo "Creating backup..."
ssh ubuntu@$JADINGIP "cp /home/ubuntu/.bashrc /home/ubuntu/app/.bashrc"
echo "..."
ssh ubuntu@$JADINGIP "tar -czf /home/ubuntu/app.tar.gz -C /home/ubuntu app"
echo "Backup created"

echo "Downloading backup..."
scp ubuntu@$JADINGIP:/home/ubuntu/app.tar.gz .
echo "Backup downloaded"

echo "Cleaning up..."
ssh ubuntu@$JADINGIP "rm /home/ubuntu/app.tar.gz"
echo "Done"
