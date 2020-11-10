package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/redsift/sandbox-go/sandbox"
)

const SANDBOX_PATH = "/run/sandbox/sift/"
const PROJECT_LOCATION = SANDBOX_PATH
const SIFT_GO_LOCATION = PROJECT_LOCATION + "/sandbox/sift.go"

const siftTemp = `package sandbox

import (
	"github.com/redsift/go-sandbox-rpc"{{range $p := .Paths}}
	"{{$p}}"{{end}}
)

var Computes = map[int]func(sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error){ {{range $i, $e := .NodeNames}}
	{{$i}} : {{$e}}.Compute,{{end}}
}`

var temp = template.Must(template.New("sift.go").Parse(siftTemp))

func main() {
	info, err := sandbox.NewInit(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(log.Lshortfile)

	uniquePaths := map[string]int{}
	nodeNames := map[int]string{}
	for _, i := range info.Nodes {
		node := info.Sift.Dag.Nodes[i]
		if node.Implementation == nil || len(node.Implementation.Go) == 0 {
			log.Fatalf("Requested to install a non-Go node at index %d", i)
		}

		implPath := node.Implementation.Go
		log.Printf("Installing node: %s : %s\n", node.Description, implPath)

		// absolutePath := path.Join(i.SIFT_ROOT, node.Implementation.Go)
		if _, err := os.Stat(path.Join(info.SIFT_ROOT, implPath)); os.IsNotExist(err) {
			log.Fatalf("Implementation at index %d : %s does not exist!", i, implPath)
		}

		packageName := path.Base(implPath)
		if strings.HasSuffix(packageName, ".go") {
			implPath = path.Dir(implPath)
			packageName = path.Base(implPath)
		}

		uniquePaths[implPath] = 1
		nodeNames[i] = packageName
	}
	paths := []string{}
	for k := range uniquePaths {
		paths = append(paths, k)
	}

	fo, err := os.Create(SIFT_GO_LOCATION)
	if err != nil {
		log.Fatal(err)
	}
	err = temp.Execute(fo, struct {
		Paths     []string
		NodeNames map[int]string
	}{paths, nodeNames})
	if err != nil {
		log.Fatalf("Failed to generate sift.go: %s", err.Error())
	}
	err = fo.Close()
	if err != nil {
		log.Fatal(err)
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
	log.Printf("%s\n", bstdoutStderr)
	if err != nil {
		log.Fatalf("Building sandbox failed: %s", err)
	}
}
