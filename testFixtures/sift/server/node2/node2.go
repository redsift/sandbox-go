package node2

import (
  "fmt"
  rpc "github.com/redsift/go-sandbox-rpc"
  "sift/utils"
)

func Compute(rpc.ComputeRequest) ([]rpc.ComputeResponse, error){
  utils.Greet()
  fmt.Println("2: helllloooo wordddd!")
  return nil, nil
}