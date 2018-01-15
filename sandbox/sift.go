package sandbox

import (
	"github.com/redsift/go-sandbox-rpc"
)

var Computes = map[int]func(sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error){}
