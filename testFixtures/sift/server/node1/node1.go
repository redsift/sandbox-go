package node1

import (
  "fmt"
  rpc "github.com/redsift/go-sandbox-rpc"
)

func Compute(rpc.ComputeRequest) ([]rpc.ComputeResponse, error){

  fmt.Println("helllloooo worldddd!")
  return nil, nil
}