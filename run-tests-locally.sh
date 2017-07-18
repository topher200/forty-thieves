#!/bin/bash

# uses the docker postgres database
export DSN="postgres://postgres@localhost:5432/forty-thieves?sslmode=disable"

go test ./... "$@"

retval=$?
if [ $retval != 0 ]; then
    echo "tests error! exit code: $retval"
    exit $retval
fi
