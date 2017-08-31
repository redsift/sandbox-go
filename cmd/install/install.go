package main

import (
	"sandbox-go/sandbox"
	"os"
	"fmt"
  "path"
  "strings"
  "strconv"
  "io/ioutil"
)

const SIFT_GO_LOCATION = "/usr/lib/redsift/sandbox/src/sandbox-go/sandbox/sift.go"
const firstPart = `package sandbox

import (
  rpc "github.com/redsift/go-sandbox-rpc"
`

const secondPart = `
)

type RedsiftFunc func(rpc.ComputeRequest) ([]rpc.ComputeResponse, error)

var Computes = map[int]RedsiftFunc{`

func main(){
	info, err := sandbox.NewInit(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

  fp := firstPart
  sp := secondPart
  for _, i := range(info.Nodes){
    node := info.Sift.Dag.Nodes[i]
    if node.Implementation == nil || len(node.Implementation.Go) == 0 {
      fmt.Printf("Requested to install a non-Go node at index %d\n",i)
      os.Exit(1)
    }

    implPath := node.Implementation.Go
    fmt.Printf("Installing node: %s : %s\n", node.Description, implPath);

    // absolutePath := path.Join(i.SIFT_ROOT, node.Implementation.Go)
    if _, err := os.Stat(path.Join(info.SIFT_ROOT, implPath)); os.IsNotExist(err) {
       fmt.Printf("Implementation at index %d : %s does not exist!\n", i, implPath)
       os.Exit(1)
    }

    packageName := path.Base(implPath)
    if strings.HasSuffix(packageName, ".go"){
      implPath = path.Dir(implPath)
      packageName = path.Base(implPath)
    }

    fp += "\n  \"" + strings.Replace(implPath, "server", "sandbox-go/sandbox/sift", 1) + "\""
    sp += "\n  " + strconv.Itoa(i) + ": " + packageName + ".Compute,"
  }
  sp += "}"

  err = ioutil.WriteFile(SIFT_GO_LOCATION, []byte(fp+sp), 0644)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

