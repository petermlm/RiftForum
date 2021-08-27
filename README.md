This is RiftForum, an implementation of a simple forum inspired by the old
phpBB and Simple Machines Forums.

It supports:

 * User groups
 * Registration by invite
 * BB code for formating
 * Bots

# Dependencies

Aside from the dependencies listed in `go.mod`, RiftForum requires PostgreSQL
and Redis to work. A docker-compose file is provided to handle both of these
dependencies.

[Air](https://github.com/cosmtrek/air) may also be used for development, but
isn't fully necessary.

# Development

## Running

Start postgres and redis using the docker compose file:

    docker-compose up

Run the migrations:

    make migrations

Then run the server:

    make run

Alternatively, [air](https://github.com/cosmtrek/air) can be used for live
reload, if it is configured:

    make air

## PSQL

To run psql, do:

    docker-compose exec postgres psql -h localhost -U riftforum_user -d riftforum_db
