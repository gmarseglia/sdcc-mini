#!/bin/bash

source ./Docker/names.config

echo $SERVER_IMAGE_NAME

docker build -t $SERVER_IMAGE_NAME -f $SERVER_IMAGE_PATH .