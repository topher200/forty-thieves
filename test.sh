#!/bin/bash

GOCACHE=off go test $(go list ./... | grep -v vendor) "$@"

retval=$?
if [ $retval != 0 ]; then
    echo "tests error! exit code: $retval"
    exit $retval
fi
