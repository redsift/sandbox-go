package main

import (
	"fmt"
	"os"
	"sandbox-go/sandbox"
	"sync"
	"time"

	ms "github.com/redsift/go-mangosock"
	s "github.com/redsift/go-socket"
	"runtime/debug"
	"errors"
)

func main() {
	info, err := sandbox.NewInit(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, i := range info.Nodes {
		node := info.Sift.Dag.Nodes[i]
		if node.Implementation == nil || len(node.Implementation.Go) == 0 {
			fmt.Printf("Requested to run a non-Go node at index %d\n", i)
			os.Exit(1)
		}

		implPath := node.Implementation.Go
		fmt.Printf("Running node: %s : %s\n", node.Description, implPath)

		if info.DRY {
			continue
		}
		wg.Add(1)
		url := fmt.Sprintf("ipc://%s/%d.sock", info.IPC_ROOT, i)
		go func(url string, idx int) {
			defer wg.Done()

			canSend := false
			var sock s.Socket
			var err error
			var msg []byte
			if sock, err = ms.NewRepSocket(); err != nil {
				die("can't get new rep socket: %s", err)
			}

			if err = sock.Connect(url); err != nil {
				die("can't dial on rep socket: %s", err)
			}

			sendErr := func(nerr error) {
				resp, err := sandbox.ToErrorBytes("error from node", nerr.Error())
				if err != nil {
					die("issue encoding your error: %s", err)
				}
				err = sock.Send(resp)
				if err != nil {
					die("can't send reply: %s", err)
				}
				canSend = false
				time.Sleep(1 * time.Second)
			}

			defer func (){
				event := recover()
				if event != nil {
					fmt.Printf("Stack: %s", debug.Stack())

					err := errors.New("panic")
					if evErr, ok := event.(error); ok {
						err = evErr
					}

					// if can send, then send err back on socket
					if canSend {
						sendErr(err)
					}
					die("caught a node panic: %s", err)
				}
			}()

			for {
				canSend = false

				msg, err = sock.Recv()
				if err != nil {
					die("error receiving from socket: %s", err)
				}
				canSend = true

				cr, err := sandbox.FromEncodedMessage(msg)
				if err != nil {
					sendErr(err)
					die("can't decode message: %s", err)
				}
				if _, ok := sandbox.Computes[idx]; !ok {
					sendErr(fmt.Errorf("no node with id: %d", idx))
					die("no node with id: %d", idx)
				}
				start := time.Now()
				nresp, err := sandbox.Computes[idx](cr)
				if err != nil {
					sendErr(err)
					continue
				}
				end := time.Since(start)
				t := []int64{int64(end / time.Second)}
				t = append(t, int64(end)-t[0])


				resp, err := sandbox.ToEncodedMessage(nresp, t)
				if err != nil {
					sendErr(err)
					die("issue encoding your response: %s", err)
				}

				err = sock.Send(resp)
				if err != nil {
					die("can't send reply: %s", err)
				}
				canSend = false
			}
		}(url, i)
	}
	wg.Wait()

	if info.DRY {
		os.Exit(0)
	}
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}
