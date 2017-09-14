#!/bin/bash

# git clone https://github.com/redsift/go-sandbox-rpc root/usr/lib/redsift/workspace/src/github.com/redsift/go-sandbox-rpc

DEV_LOC=/Users/chrisvon/Documents/Developer/redsift
docker run \
-v ${DEV_LOC}/sandbox-go/testFixtures/sift:/run/sandbox/sift \
-ti gotest

# -v ${DEV_LOC}/sandbox-go:/usr/lib/redsift/sandbox/src/sandbox-go \
