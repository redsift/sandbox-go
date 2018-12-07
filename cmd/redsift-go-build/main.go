package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/redsift/sandbox-go/cmd/redsift-go-build/internal"

	"github.com/redsift/go-sandbox-rpc/sift"
)

const MainName = "main___.go"

var deferFunc = func() {}

func main() {
	siftRoot, nodesPackage, nodes, err := configure(os.Args[1:])
	if err != nil {
		die("invalid sandbox configuration: %s", err)
	}

	uniquePaths := make(map[string]struct{})
	for _, n := range nodes {
		uniquePaths[n.PkgPath] = struct{}{}
		log.Printf("selecting node #%d %q", n.Idx, n.Desc)
	}

	mainAbsPath := path.Join(siftRoot, nodesPackage, MainName)
	log.Printf("generating %q", mainAbsPath)
	cfg := struct {
		Paths map[string]struct{}
		Nodes []Node
	}{
		Paths: uniquePaths,
		Nodes: nodes,
	}
	//deferFunc = func() {
	//	log.Printf("deleting %s", mainAbsPath)
	//	_ = os.Remove(mainAbsPath)
	//}
	defer deferFunc()
	if err := internal.GenerateSiftMain(MainName, mainAbsPath, cfg); err != nil {
		die("couldn't generate %q: %s", mainAbsPath, err)
	}

	args := []string{"build"}
	if os.Getenv("LOG_LEVEL") == "debug" {
		args = append(args, "-x")
	}
	args = append(args, "-v", "-o", envString("SIFT_BIN", "/run/sandbox/sift/server/_run"), MainName)
	buildCmd := exec.Command("go", args...)
	//buildCmd := exec.Command("go", "env")
	buildCmd.Env = []string{fmt.Sprintf("GOPATH=%s", os.Getenv("GOPATH"))}
	// TODO
	//  create link $GOPATH/$nodesPackage
	//  remove link $GOPATH/$nodesPackage if we created it
	buildCmd.Dir = path.Join(os.Getenv("GOPATH"), "src", nodesPackage)
	log.Printf("building %q (dir=%q, env=%v, args=%v)", mainAbsPath, buildCmd.Dir, buildCmd.Env, args)

	r, w := io.Pipe()
	buildCmd.Stdout = w
	buildCmd.Stderr = w

	if err := buildCmd.Start(); err != nil {
		die("can't start 'go build': %s", err)
	}

	go io.Copy(os.Stdout, r)

	if err := buildCmd.Wait(); err != nil {
		die("'go build' failed: %s", err)
	}
}

func envString(key, def string) string {
	v, found := os.LookupEnv(key)
	if !found {
		return def
	}
	return v
}

func die(format string, v ...interface{}) {
	deferFunc()
	log.Println("FATAL:", fmt.Sprintf(format, v...))
	os.Exit(1)
}

type Node struct {
	Idx     int
	Desc    string
	PkgPath string
	PkgName string
}

func configure(args []string) (string, string, []Node, error) {
	newError := func(f string, args ...interface{}) (string, string, []Node, error) {
		return "", "", nil, fmt.Errorf(f, args...)
	}
	errEnvVarNotFound := func(s string) (string, string, []Node, error) {
		return newError("environment variable %s not found", s)
	}

	if len(args) == 0 {
		return newError("no nodes requested; nothing to do")
	}

	var (
		siftRoot string
		siftFile string
		found    bool
	)
	if siftRoot, found = os.LookupEnv("SIFT_ROOT"); !found {
		return errEnvVarNotFound("SIFT_ROOT")
	}
	if siftFile, found = os.LookupEnv("SIFT_JSON"); !found {
		return errEnvVarNotFound("SIFT_ROOT")
	}

	siftPath := path.Join(siftRoot, siftFile)
	raw, err := ioutil.ReadFile(siftPath)
	if err != nil {
		return newError("could't read file %s: %s", siftPath, err)
	}

	var root sift.Root

	if err := json.Unmarshal(raw, &root); err != nil {
		return newError("couldn't parse %s: %s", siftPath, err)
	}

	if !root.HasDag() {
		return newError("%s has not any nodes", siftPath)
	}

	requiredNodes := make(map[int]struct{})
	for i, arg := range args {
		n, err := strconv.Atoi(arg)
		if err != nil {
			return newError("can't parse argument #%d %q: %s", i, arg, err)
		}
		requiredNodes[n] = struct{}{}
	}

	nodesPackage := ""
	nodeRoots := make(map[string]struct{})
	var nodes []Node
	for i, n := range root.Dag.Nodes {
		if _, found := requiredNodes[i]; !found {
			continue
		}
		if n.Implementation == nil || len(n.Implementation.Go) == 0 {
			return newError("requested Node #%d %q is non-Go Node", i, n.Description)
		}
		nodePath := n.Implementation.Go
		if file := path.Base(nodePath); !strings.HasSuffix(file, ".go") {
			return newError("requested Node #%d %q implementation is non-Go: %s", i, n.Description, file)
		}
		if _, err := os.Stat(path.Join(siftRoot, nodePath)); os.IsNotExist(err) {
			return newError("requested Node #%d %q has not implemented: %s", i, n.Description, err)
		}
		// nodePath := "server/package/Node.go"
		pkgPath := path.Dir(nodePath) // "server/package"
		pkgName := path.Base(pkgPath) // "package"
		nodes = append(nodes, Node{
			Idx:     i,
			Desc:    n.Description,
			PkgPath: pkgPath,
			PkgName: pkgName,
		})

		nodesPackage = nodePath[:strings.Index(nodePath, "/")+1]
		nodeRoots[nodesPackage] = struct{}{}

		delete(requiredNodes, i)
	}

	if len(nodeRoots) != 1 {
		return newError("no common root folder found for requested implementations")
	}

	// have all nodes been resolved ?
	if len(requiredNodes) != 0 {
		n := make([]int, 0, len(requiredNodes))
		for i := range requiredNodes {
			n = append(n, i)
		}
		return newError("requested to install undefined nodes %v", n)
	}

	return siftRoot, nodesPackage, nodes, nil
}
