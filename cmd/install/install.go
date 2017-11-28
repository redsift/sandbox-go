package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sandbox-go/sandbox"
	"strings"
	"text/template"
	"runtime/debug"
)

const SANDBOX_PATH = "sandbox-go"
const PROJECT_LOCATION = "/usr/lib/redsift/workspace/src/" + SANDBOX_PATH
const SIFT_GO_LOCATION = PROJECT_LOCATION + "/sandbox/sift.go"
const sift_temp = `package sandbox

import (
	"github.com/redsift/go-sandbox-rpc"{{range $i, $e := .}}
	"{{$e.Path}}"{{end}}
)

var Computes = map[int]func(sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error){ {{range $i, $e := .}}
	{{$i}} : {{$e.Name}}.Compute,{{end}}
}`

func main() {
	info, err := sandbox.NewInit(os.Args[1:])
	if err != nil {
		die("%s", err.Error())
	}

	debug.SetMaxThreads(10)

	type Anode struct {
		Path string
		Name string
	}
	availableNodes := map[int]Anode{}
	for _, i := range info.Nodes {
		node := info.Sift.Dag.Nodes[i]
		if node.Implementation == nil || len(node.Implementation.Go) == 0 {
			die("Requested to install a non-Go node at index %d\n", i)
		}

		implPath := node.Implementation.Go
		fmt.Printf("Installing node: %s : %s\n", node.Description, implPath)

		// absolutePath := path.Join(i.SIFT_ROOT, node.Implementation.Go)
		if _, err := os.Stat(path.Join(info.SIFT_ROOT, implPath)); os.IsNotExist(err) {
			die("Implementation at index %d : %s does not exist!\n", i, implPath)
		}

		packageName := path.Base(implPath)
		if strings.HasSuffix(packageName, ".go") {
			implPath = path.Dir(implPath)
			packageName = path.Base(implPath)
		}
		availableNodes[i] = Anode{
			Path: implPath,
			Name: packageName,
		}
	}

	fo, err := os.Create(SIFT_GO_LOCATION)
	if err != nil {
		die("%s", err.Error())
	}
	defer func() {
		if err := fo.Close(); err != nil {
			die("%s", err.Error())
		}
	}()

	t := template.New("sift.go")
	t, _ = t.Parse(sift_temp)
	err = t.Execute(fo, availableNodes)
	if err != nil {
		die("Failed to generate sift.go: %s", err.Error())
	}

	//
	// Build Phase
	//
	buildArgs := []string{"build"}
	if os.Getenv("LOG_LEVEL") == "debug" {
		buildArgs = append(buildArgs, "-x")
	}
	buildArgs = append(buildArgs, "-v", "-o", "/run/sandbox/sift/server/_run", path.Join(PROJECT_LOCATION, "cmd/run/run.go"))
	bcmd := exec.Command("go", buildArgs...)
	bstdoutStderr, err := bcmd.CombinedOutput()
	fmt.Printf("%s\n", bstdoutStderr)
	if err != nil {
		die("Building sandbox failed: %s", err)
	}
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}
