#!/bin/bash

go test ./... "$@"

retval=$?
if [ $retval != 0 ]; then
    echo "tests error! exit code: $retval"
    exit $retval
fi
