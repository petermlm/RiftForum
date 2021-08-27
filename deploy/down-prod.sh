#!/bin/bash

cd docker
docker-compose -p riftforum -f docker-compose-prod.yml down
