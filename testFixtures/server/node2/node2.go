package node2

import (
  "fmt"
  rpc "github.com/redsift/go-sandbox-rpc"
)

func Compute(rpc.ComputeRequest) ([]rpc.ComputeResponse, error){

  fmt.Println("2: helllloooo wordddd!")
  return nil, nil
}