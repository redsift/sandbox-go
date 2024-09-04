package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/redsift/sandbox-go/sandbox"
	"github.com/redsift/sandbox-go/sandbox/modedit"
)

const PROJECT_LOCATION = "/build/"

const SIFT_GO_LOCATION = PROJECT_LOCATION + "sandbox/sift.go"

const siftTemp = `
package sandbox

import (
	"github.com/redsift/go-sandbox-rpc"{{range $p := .Paths}}
	"{{$p}}"{{end}}
)

var Computes = map[int]func(sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error){ {{range $i, $e := .NodeNames}}
	{{$i}} : {{$e}}.Compute,{{end}}
}`

var siftTemplate = template.Must(template.New("sift.go").Parse(siftTemp))

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
	err = siftTemplate.Execute(fo, struct {
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
	// copy replace directives from sift go.mod to sandbox go.mod

	sbxMod := filepath.Join(PROJECT_LOCATION, "go.mod")
	modedit.CopyReplace(filepath.Join(info.SIFT_ROOT, "server", "go.mod"), sbxMod, sbxMod)

	sbxSum := filepath.Join(PROJECT_LOCATION, "go.sum")
	modedit.CopySum(filepath.Join(info.SIFT_ROOT, "server", "go.sum"), sbxSum, sbxSum)

	mode := "mod"
	vendor := filepath.Join(PROJECT_LOCATION, "vendor")
	if s, err := os.Stat(vendor); err == nil && s.IsDir() {
		mode = "vendor"
		run("go", "mod", "vendor")
	}

	buildArgs := []string{"build", "-mod", mode}
	if os.Getenv("LOG_LEVEL") == "debug" {
		buildArgs = append(buildArgs, "-x")
	}

	buildArgs = append(buildArgs,
		"-v",
		"-o", info.Output,
		path.Join(PROJECT_LOCATION, "cmd/run/run.go"),
	)

	run("go", buildArgs...)

	log.Printf("Installed nodes: %v : %v", nodeNames, uniquePaths)
}

func run(c string, args ...string) {
	bcmd := exec.Command(c, args...)
	bcmd.Stdout = os.Stdout
	bcmd.Stderr = os.Stderr
	bcmd.Dir = PROJECT_LOCATION

	if err := bcmd.Run(); err != nil {
		log.Fatalf("Building sandbox failed: %s", err)
	}
}
