#!/bin/bash

DEV_LOC=/Users/chrisvon/Documents/Developer/redsift
docker run \
-v ${DEV_LOC}/sandbox-go:/tmp/sandbox \
-v ${DEV_LOC}/sifts/hello-go-sift:/run/sandbox/sift \
-ti gotest

# -v ${DEV_LOC}/sandbox-go:/usr/lib/redsift/sandbox/src/sandbox-go \
# -v ${DEV_LOC}/sandbox-swift/TestFixtures/sift:/run/sandbox/sift \
# -e SIFT_ROOT=/tmp/sandbox/sift \