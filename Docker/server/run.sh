#!/bin/bash

source ./Docker/names.config

if [[ $# -eq 1 ]] && [[ $1 -eq "-i" ]]; then
    docker run --rm -it -p 55555:55555/tcp -p 55556:55556/tcp --name=$SERVER_CONTAINER_NAME $SERVER_IMAGE_NAME:latest
else
    docker run --rm -d  -p 55555:55555/tcp -p 55556:55556/tcp --name=$SERVER_CONTAINER_NAME $SERVER_IMAGE_NAME:latest
fi

