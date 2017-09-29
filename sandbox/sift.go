package sandbox

import (
	"github.com/redsift/go-sandbox-rpc"
)

// type RedsiftFunc func(sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error)

var Computes = map[int]func(sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error){}
