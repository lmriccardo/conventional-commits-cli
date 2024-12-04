#!/bin/bash

dockerfile=$1

if [ -z $dockerfile ];
then
    echo "A dockerfile is required"
    exit
fi

docker build -f $dockerfile -t ccommits-cli:latest .
docker tag ccommits-cli:latest lmriccardo/ccommits-cli:latest
docker login
docker push lmriccardo/ccommits-cli:latest