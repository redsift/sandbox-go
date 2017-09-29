package node2

import (
	"fmt"
	"sift/utils"

	"github.com/redsift/go-sandbox-rpc"
)

func Compute(sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error) {
	utils.Greet()
	fmt.Println("2: helllloooo wordddd!")

	resp := []sandboxrpc.ComputeResponse{
		sandboxrpc.NewComputeResponse("api_rpc", "1341231", []byte("first payload"), 0, 0),
		sandboxrpc.NewComputeResponse("stats", "index_stats", []byte("second payload"), 0, 0),
	}
	return resp, nil
}
