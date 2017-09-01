package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"

	ms "github.com/redsift/go-mangosock"
	rpc "github.com/redsift/go-sandbox-rpc"
	s "github.com/redsift/go-socket"
)

func TestComputeRequest(t *testing.T) {
	var sock s.Socket
	var err error
	var msg []byte

	url := "ipc:///run/sandbox/ipc/1.sock"
	t.Log("will send request")
	if sock, err = ms.NewReqSocket(); err != nil {
		t.Logf("can't get new req socket: %s", err.Error())
	}
	if err = sock.Bind(url); err != nil {
		t.Logf("can't dial on req socket: %s", err.Error())
	}

	b, err := json.Marshal(rpc.ComputeRequest{})
	if err = sock.Send(b); err != nil {
		t.Logf("can't send message on push socket: %s", err.Error())
	}
	if msg, err = sock.Recv(); err != nil {
		t.Logf("can't receive date: %s", err.Error())
	}
	t.Logf("and this is what I got back %s", msg)
	sock.Close()
}

func TestMain(m *testing.M) {
	go func() {
		cmd := exec.Command("go", "run", "/usr/lib/redsift/sandbox/src/sandbox-go/cmd/run/run.go", "1")
		stdoutStderr, _ := cmd.CombinedOutput()
		fmt.Printf("%s\n", stdoutStderr)
	}()
	os.Exit(m.Run())
}
