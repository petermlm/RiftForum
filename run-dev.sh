#!/bin/bash

export UID=$(id -u)
export GID=$(id -g)
cd docker
docker-compose -f docker-compose-dev.yml up $@
