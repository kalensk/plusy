#!/bin/bash

docker run \
    --rm \
    -p 7473:7473 \
    -p 7474:7474 \
    -p 7687:7687 \
    --volume=$HOME/neo4j/data:/data \
    --volume=$HOME/neo4j/logs:/logs \
    --name neo4j \
    neo4j:3.5
