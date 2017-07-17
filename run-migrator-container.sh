#!/bin/bash

docker build -f Dockerfile_migrator -t fortythieves_migrator .
docker run -it --network fortythieves_default fortythieves_migrator /bin/bash
