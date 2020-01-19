#!/bin/bash

docker build -t livecam-dev . -f Dockerfile.dev
docker run --name=livecam --rm -it -v $(pwd)/config:/etc/livecam -v $(pwd):/app -p 3000:3000 livecam-dev
docker rmi livecam-dev
