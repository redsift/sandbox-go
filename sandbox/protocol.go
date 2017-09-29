package sandbox

import (
	"encoding/json"

	"github.com/redsift/go-sandbox-rpc"
)

func ToEncodedMessage(data []sandboxrpc.ComputeResponse, diff []int64) ([]byte, error) {
	var pd []*sandboxrpc.ComputeResponse
	for _, d := range data {
		pd = append(pd, &d)
	}
	return json.Marshal(sandboxrpc.Response{
		Out:   pd,
		Stats: map[string][]int64{"results": diff},
	})
}

func ToErrorBytes(message string, stack string) ([]byte, error) {
	return json.Marshal(sandboxrpc.Response{
		Error: map[string]string{"message": message, "stack": stack}})
}

func FromEncodedMessage(bytes []byte) (sandboxrpc.ComputeRequest, error) {
	cr := sandboxrpc.ComputeRequest{}
	err := json.Unmarshal(bytes, &cr)
	return cr, err
}
