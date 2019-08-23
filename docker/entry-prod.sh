#!/bin/bash

# Just wait five seconds before Postgres is ready
sleep 5

/go/tmp/server -migrate
/go/tmp/server
