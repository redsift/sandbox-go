package sandbox

import (
	"context"

	sandboxrpc "github.com/redsift/go-sandbox-rpc"
)

var Computes = map[int]func(context.Context, sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error){}
