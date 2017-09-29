#!/bin/bash
rm -rf ./sift ./sift.go

cat > sift.go << EOF
package sandbox

import (
	"github.com/redsift/go-sandbox-rpc"
)

var Computes = map[int]func(sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error){}
EOF