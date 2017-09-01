package main

import (
	"fmt"
	"os"
	"sandbox-go/sandbox"
	"sync"
	"time"

	ms "github.com/redsift/go-mangosock"
	s "github.com/redsift/go-socket"
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
		go func(url string) {
			defer wg.Done()
			var sock s.Socket
			var err error
			var msg []byte
			if sock, err = ms.NewRepSocket(); err != nil {
				die("can't get new rep socket: %s", err)
			}

			if err = sock.Connect(url); err != nil {
				die("can't dial on rep socket: %s", err.Error())
			}
			for {
				msg, err = sock.Recv()
				cr, err := sandbox.FromEncodedMessage(msg)
				if err != nil {
					die("can't decode message: %s", err.Error())
				}

				start := time.Now()
				nresp, nerr := sandbox.Computes[i](cr)
				end := time.Since(start)
				t := []int64{int64(end / time.Second)}
				t = append(t, int64(end)-t[0])

				var resp []byte
				if nerr == nil {
					resp, err = sandbox.ToEncodedMessage(nresp, t)
					if err != nil {
						die("issue encoding your response: %s", err.Error())
					}
				} else {
					resp, err = sandbox.ToErrorBytes("error from node", nerr.Error())
					if err != nil {
						die("issue encoding your error: %s", err.Error())
					}
				}

				err = sock.Send(resp)
				if err != nil {
					die("can't send reply: %s", err.Error())
				}
			}
		}(url)
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
