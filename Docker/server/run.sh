#!/bin/bash

source ./Docker/names.config

if [[ $# -eq 1 ]] && [[ $1 -eq "-i" ]]; then
    docker run --rm -it --network host --name=$SERVER_CONTAINER_NAME $SERVER_IMAGE_NAME:latest
else
    docker run --rm -d --network host --name=$SERVER_CONTAINER_NAME $SERVER_IMAGE_NAME:latest
fi

