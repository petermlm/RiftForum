#!/bin/bash

# Just wait five seconds before Postgres and Redis are ready
sleep 5

/bin/riftforum -migrate
/bin/riftforum
