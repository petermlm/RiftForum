#!/bin/bash

# Just wait five seconds before Postgres and Redis are ready
sleep 5

/go/tmp/server -migrate
/go/tmp/server
