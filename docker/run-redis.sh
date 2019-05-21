#!/bin/bash

docker run \
    --rm \
    -p 6379:6379 \
    --mount type=volume,source=redis-data,target=/data \
    --name redis \
    redis:5.0-rc-stretch
