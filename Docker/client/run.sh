#!/bin/bash

source ./Docker/names.config

if [[ $# -eq 1 ]]; then

    if [[ $1 -eq "-i" ]]; then
        docker run --rm -it --name=$CLIENT_CONTAINER_NAME $CLIENT_IMAGE_NAME:latest
    fi

    if [[ $1 -eq "-e" ]]; then
        docker run --rm -d -e FrontAddr='192.168.1.49' --name=$CLIENT_CONTAINER_NAME $CLIENT_IMAGE_NAME:latest
    fi

else
    docker run --rm -d --name=$CLIENT_CONTAINER_NAME $CLIENT_IMAGE_NAME:latest
fi

