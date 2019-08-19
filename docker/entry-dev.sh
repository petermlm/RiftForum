#!/bin/bash

go run src/*.go -migrate
/air_run/air -c air.conf
