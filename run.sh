#!/bin/bash

# Exit on error
set -e

# install then run
go build
./forty-thieves "$@"
