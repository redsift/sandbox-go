#!/bin/bash
rm -rf ./sift ./sift.go

cat > sift.go << EOF
package sandbox

import (
	rpc "github.com/redsift/go-sandbox-rpc"
)

var Computes = map[int]func(rpc.ComputeRequest) ([]rpc.ComputeResponse, error){}
EOF