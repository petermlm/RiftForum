This is RiftForum, an implementation of a simple forum inspired by the old
phpBB and Simple Machines Forums.

It supports:

 * User groups
 * Registration by invite
 * BB code for formatting
 * Bots

RiftForum's Frontend is developed using Go templates and Boostrap. It's
JavaScript free otherwise!

# Bots

There are a few bots in RiftForum. The bots see every new user, topic, and
message and may act upon them. The bots which are currently developed are:

 * `GreeterBot` - Creates a new topic for every new user.
 * `RedditBot` - Every time someone mentions a subreddit, this bot creates a
                 message with the current top 5 posts of the hot tab of that
                 subreddit.
 * `YoutubeBot` - This bot acts if someone writes `!youtubelist`. It creates a
                  message with a list of every youtube video on that topic.

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
