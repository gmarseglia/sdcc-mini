#!/bin/bash

source ./Docker/names.config

re='^[0-9]+$'

if [[ $# -eq 1 ]] && [[ $1 =~ $re ]]; then
    for i in $(seq 1 1 $1); do
        docker run --rm -d -p 55557/tcp $WORKER_IMAGE_NAME:latest
    done
else
    docker run --rm -d -p 55557/tcp $WORKER_IMAGE_NAME:latest
fi

