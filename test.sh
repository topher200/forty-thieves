#!/bin/bash

go test ./... -v

retval=$?
if [ $retval != 0 ]; then
    echo "tests error! exit code: $retval"
fi
