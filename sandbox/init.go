package sandbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

  rpc "github.com/redsift/go-sandbox-rpc"
	"github.com/redsift/go-sandbox-rpc/sift"
)

type RedsiftFunc func(request rpc.ComputeRequest) (rpc.ComputeResponse, error)

type Init struct {
	SIFT_ROOT string
	SIFT_JSON string
	IPC_ROOT  string
	DRY       bool
	sift      sift.Root
	nodes     []int
}

func NewInit(args []string) (Init, error) {
	if args == nil || len(args) == 0 {
		return Init{}, errors.New("No nodes to execute")
	}
	i := Init{
		SIFT_ROOT: os.Getenv("SIFT_ROOT"),
		SIFT_JSON: os.Getenv("SIFT_JSON"),
		IPC_ROOT:  os.Getenv("IPC_ROOT"),
		DRY:       false,
		nodes:     []int{},
	}

	if len(i.SIFT_ROOT) == 0 {
		return Init{}, errors.New("Environment SIFT_ROOT not set")
	}

	if len(i.SIFT_JSON) == 0 {
		return Init{}, errors.New("Environment SIFT_JSON not set")
	}

	if len(i.IPC_ROOT) == 0 {
		return Init{}, errors.New("Environment IPC_ROOT not set")
	}

	if len(os.Getenv("DRY")) > 0 {
		fmt.Println("Unit Test Mode")
		i.DRY = true
	}

	raw, err := ioutil.ReadFile(path.Join(i.SIFT_ROOT, i.SIFT_JSON))
	if err != nil {
		return Init{}, err
	}

	err = json.Unmarshal(raw, &i.sift)
	if err != nil {
		return Init{}, err
	}

	if !i.sift.HasDag() {
		return Init{}, errors.New("sift.json does not contain any nodes")
	}

	for _, v := range args {
		a, _ := strconv.Atoi(v)
		i.nodes = append(i.nodes, a)
	}
	return i, nil
}
