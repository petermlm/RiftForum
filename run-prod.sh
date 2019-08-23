#!/bin/bash

export UID=$(id -u)
export GID=$(id -g)
cd docker
docker-compose -d -f docker-compose-prod.yml $@
