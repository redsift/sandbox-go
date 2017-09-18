package sandbox

import (
	"encoding/json"

	rpc "github.com/redsift/go-sandbox-rpc"
)

func ToEncodedMessage(data []rpc.ComputeResponse, diff []int64) ([]byte, error) {
	var pd []*rpc.ComputeResponse
	for _, d := range(data){
		pd = append(pd, &d)
	}
	return json.Marshal(rpc.Response{
		Out:   pd,
		Stats: map[string][]int64{"results": diff},
	})
}

func ToErrorBytes(message string, stack string) ([]byte, error) {
	return json.Marshal(rpc.Response{
		Error: map[string]string{"message": message, "stack": stack}})
}

func FromEncodedMessage(bytes []byte) (rpc.ComputeRequest, error) {
	cr := rpc.ComputeRequest{}
	err := json.Unmarshal(bytes, &cr)
	return cr, err
}
