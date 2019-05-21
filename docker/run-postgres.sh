#!/bin/bash

docker run \
    --rm \
    -p 5432:5432 \
    --mount type=volume,source=pgdata,target=/data \
    --name postgres \
    postgres:11.2
