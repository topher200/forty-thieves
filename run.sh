#!/bin/bash

# Exit on error
set -e

# Test, install, then run
go test ./...
go build
./forty-thieves --deterministic
