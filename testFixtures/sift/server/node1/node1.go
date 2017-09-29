package node1

import (
	"fmt"

	"github.com/redsift/go-sandbox-rpc"
)

func Compute(sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error) {

	fmt.Println("helllloooo worldddd!")
	return nil, nil
}
