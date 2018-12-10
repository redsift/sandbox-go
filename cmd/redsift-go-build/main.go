package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/redsift/go-sandbox-rpc/sift"
	"github.com/redsift/sandbox-go/cmd/redsift-go-build/internal"
)

var VerboseMode bool

func main() {
	var (
		siftPath     string
		reuseTempDir bool
		binPath      string
		helpWanted   bool
		allNodes     bool
	)
	flag.BoolVar(&helpWanted, "h", false, "show this help")
	flag.StringVar(&siftPath, "sift", "", "path to the sift.json file; you can set SIFT_ROOT and SIFT_JSON env vars as an alternative")
	flag.StringVar(&binPath, "o", EnvString("SIFT_BIN", "/run/sandbox/sift/server/_run"), "write the resulting executable to the named output file; you can set SIFT_BIN env var as an alternative")
	flag.BoolVar(&reuseTempDir, "reuse-workdir", false, "reuse temporary working dir")
	flag.BoolVar(&allNodes, "all", false, "compile all go nodes")
	flag.BoolVar(&VerboseMode, "v", false, "verbose mode")
	flag.Parse()

	if helpWanted {
		flag.PrintDefaults()
		os.Exit(1)
	}

	const (
		siftMainFile = "main___.go"
	)
	defer RunBeforeExit()

	siftRoot, nodesPackage, nodes, err := Configure(siftPath, flag.Args(), allNodes)
	if err != nil {
		Fatalf("invalid sandbox configuration: %s", err)
	}

	Verbosef("using sift: dir=%q", siftRoot)
	uniquePaths := make(map[string]struct{})
	for _, n := range nodes {
		uniquePaths[n.PkgPath] = struct{}{}
		Verbosef("node added: idx=%d desc=%q", n.Idx, n.Desc)
	}

	siftMainPath := path.Join(siftRoot, nodesPackage, siftMainFile)
	cfg := struct {
		Paths map[string]struct{}
		Nodes []Node
	}{
		Paths: uniquePaths,
		Nodes: nodes,
	}

	Deferf(func() { _ = os.RemoveAll(siftMainPath) }, "remove %q", siftMainPath)
	if err := internal.GenerateSiftMain(siftMainFile, siftMainPath, cfg); err != nil {
		Fatalf("couldn't create %q: %s", siftMainPath, err)
	}

	localGoPath, workingDir, err := MkTempPkgDir(path.Join(siftRoot, nodesPackage), reuseTempDir)
	if err != nil {
		Fatalf("couldn't create temp dir: %s", err)
	}

	Deferf(func() { _ = os.RemoveAll(siftMainPath) }, "remove %q", localGoPath)

	args := []string{"build"}
	if os.Getenv("LOG_LEVEL") == "debug" {
		args = append(args, "-x")
	}
	if VerboseMode {
		args = append(args, "-v")
	}

	if binPath, err = filepath.Abs(binPath); err != nil {
		Fatalf("could't output result to %q: %q", binPath, err)
	}

	args = append(args, "-o", binPath, siftMainFile)

	buildCmd := exec.Command("go", args...)

	buildCmd.Env = []string{
		"PWD=" + workingDir,                             // need that as syscall.Chdir changes working dir to real dir, not to symlink
		"GOPATH=" + localGoPath + ":" + Goenv("GOPATH"), // put local gopath in front of system wide
		"PATH=/usr/bin",                                 // go might need clang
		"GOCACHE=" + Goenv("GOCACHE"),                   // use system wide cache
		"GOARCH=" + os.Getenv("GOARCH"),
		"GOOS=" + os.Getenv("GOOS"),
	}
	buildCmd.Dir = workingDir
	Verbosef("execute go: dir=%q env=%v args=%v", buildCmd.Dir, buildCmd.Env, args)

	r, w := io.Pipe()
	buildCmd.Stdout = w
	buildCmd.Stderr = w

	if err := buildCmd.Start(); err != nil {
		Fatalf("can't start 'go build': %s", err)
	}

	go io.Copy(os.Stdout, r)

	if err := buildCmd.Wait(); err != nil {
		Fatalf("'go build' failed: %s", err)
	} else {
		Verbosef("done: output=%q", binPath)
	}
}

func Goenv(name string) string {
	var out bytes.Buffer
	cmd := exec.Command("go", "env", name)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return ""
	}
	return string(bytes.TrimSpace(out.Bytes()))
}

func MkTempPkgDir(pkgDir string, reuseDir bool) (string, string, error) {
	const parentPkgDir = "___redsift-go-build"
	var (
		gopath string
		err    error
	)
	if reuseDir {
		gopath = path.Join(os.TempDir(), parentPkgDir)
	} else {
		gopath, err = ioutil.TempDir(os.TempDir(), parentPkgDir)
	}
	if err != nil {
		return "", "", err
	}
	src := path.Join(gopath, "src")
	if err := os.MkdirAll(src, 0755); err != nil {
		_ = os.RemoveAll(gopath)
		return "", "", err
	}
	pkg := path.Join(src, path.Base(pkgDir))
	if err := os.Symlink(pkgDir, pkg); err != nil && !(reuseDir && os.IsExist(err)) {
		_ = os.RemoveAll(gopath)
		return "", "", err
	}
	return gopath, pkg, nil
}

func EnvString(key, def string) string {
	v, found := os.LookupEnv(key)
	if !found {
		return def
	}
	return v
}

type deferredFunc struct {
	f      func()
	format string
	args   []interface{}
}

var deferredFuncs []deferredFunc

func Deferf(f func(), fmt string, args ...interface{}) func() {
	deferredFuncs = append(deferredFuncs, deferredFunc{f, fmt, args})
	return f
}

func RunBeforeExit() {
	for _, v := range deferredFuncs {
		Verbosef(v.format, v.args...)
		v.f()
	}
}

func Verbosef(format string, args ...interface{}) {
	if !VerboseMode {
		return
	}
	log.Println("INFO:", fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...interface{}) {
	RunBeforeExit()
	log.Println("FATAL:", fmt.Sprintf(format, args...))
	os.Exit(1)
}

type Node struct {
	Idx     int
	Desc    string
	PkgPath string
	PkgName string
}

func Configure(siftPath string, args []string, allNodes bool) (string, string, []Node, error) {
	newError := func(f string, args ...interface{}) (string, string, []Node, error) {
		return "", "", nil, fmt.Errorf(f, args...)
	}
	errEnvVarNotFound := func(s string) (string, string, []Node, error) {
		return newError("environment variable %s not found", s)
	}

	if len(args) == 0 && !allNodes {
		return newError("no nodes requested; nothing to do")
	}

	var (
		siftRoot string
		siftFile string
		found    bool
	)

	if siftPath == "" {
		if siftRoot, found = os.LookupEnv("SIFT_ROOT"); !found {
			return errEnvVarNotFound("SIFT_ROOT")
		}
		siftFile = EnvString("SIFT_JSON", "sift.json")

		siftPath = path.Join(siftRoot, siftFile)
	} else {
		p, err := filepath.Abs(siftPath)
		if err != nil {
			return newError("couldn't resolve path %s: %s", siftPath, err)
		}
		siftRoot, siftFile = path.Split(p)
	}

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
	if len(args) == 0 && allNodes {
		for i, n := range root.Dag.Nodes {
			if n.Implementation == nil || len(n.Implementation.Go) == 0 {
				continue
			}
			requiredNodes[i] = struct{}{}
		}
	} else {
		for i, arg := range args {
			n, err := strconv.Atoi(arg)
			if err != nil {
				return newError("can't parse argument #%d %q: %s", i, arg, err)
			}
			requiredNodes[n] = struct{}{}
		}
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

	// do we have all nodes resolved?
	if len(requiredNodes) != 0 {
		n := make([]int, 0, len(requiredNodes))
		for i := range requiredNodes {
			n = append(n, i)
		}
		return newError("requested to install undefined nodes %v", n)
	}

	if len(nodeRoots) != 1 {
		return newError("no common root folder found for requested implementations")
	}

	return siftRoot, nodesPackage, nodes, nil
}
