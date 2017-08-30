package sandbox

import (
  "encoding/json"
  rpc "github.com/redsift/go-sandbox-rpc"
)

func toEncodedMessage(data []*rpc.ComputeResponse, diff []float64) ([]byte, error) {
	return json.Marshal(rpc.Response{
		Out:   data,
		Stats: map[string][]float64{"results": diff},
	})
}

func toErrorBytes(message string, stack string) ([]byte, error) {
	return json.Marshal(rpc.Response{
		Error: map[string]string{"message": message, "stack": stack}})
}

func fromEncodedMessage(bytes []byte) (rpc.ComputeRequest, error) {
	cr := rpc.ComputeRequest{}
	err := json.Unmarshal(bytes, &cr)
	return cr, err
}
