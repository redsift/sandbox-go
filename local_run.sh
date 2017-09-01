#!/bin/bash

DEV_LOC=/Users/chrisvon/Documents/Developer/redsift
docker run \
-v ${DEV_LOC}/sifts/hello-go-sift:/run/sandbox/sift \
-ti gotest