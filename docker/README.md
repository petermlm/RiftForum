# Docker Files

There is a docker file for the server container and for the postgres container.

## Server

The docker file for the server Ubuntu 18.04 LTS. It installs the dependencies
to run the golang server for RiftForum's backend.

The `/app` directory contains the whole repository.

## Postgres

This docker is based on postgres 9.6 and defines three things:

 * Database name
 * Database user
 * Database pass

# Docker Compose

The docker compose files are used to build and run the containers. There are
two files:

 * `docker-compose-prod.sh`
 * `docker-compose-dev.sh`

## Production

The production file is to be used in production so the database is persistent.
The mapped directory with the postgres data is:

    REPO/docker/data

## Development

The development file is to be used while developing. The database is not
persistent.

## Build

Use the following commands to build:

    docker-compose -f docker-compose-prod.yml build
    docker-compose -f docker-compose-dev.yml build

To run:

    docker-compose -f docker-compose-prod.yml up
    docker-compose -f docker-compose-dev.yml up

To tear down:

    docker-compose -f docker-compose-prod.yml down
    docker-compose -f docker-compose-dev.yml down

## Bash and PSQL

To open a shell in the node docker use:

    docker exec -it docker_server_1 bash

This will allow you to create models, seeds, and run migrations when in
development.

To open psql, first run bash into the Postgres container, then run:

    psql -h postgres -U riftforum_user -d riftforum_db

You will be prompted for the password.
