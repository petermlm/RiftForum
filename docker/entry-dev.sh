#!/bin/bash

# Just wait five seconds before Postgres is ready
sleep 5

# Start go
cd src
go run *.go
