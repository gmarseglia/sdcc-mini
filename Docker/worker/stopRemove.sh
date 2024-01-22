#!/bin/bash

source ./Docker/names.config

 docker stop $(docker ps -a --filter ancestor=$WORKER_IMAGE_NAME --format="{{.ID}}")
