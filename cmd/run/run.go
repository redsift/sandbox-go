package main

import (
	"fmt"
	"os"
	"sandbox-go/sandbox"
	"sync"
	"time"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/rep"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
)

func main() {
	info, err := sandbox.NewInit(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	for _, i := range info.Nodes {
		node := info.Sift.Dag.Nodes[i]
		if node.Implementation == nil || len(node.Implementation.Go) == 0 {
			fmt.Printf("Requested to run a non-Go node at index %d\n", i)
			os.Exit(1)
		}

		implPath := node.Implementation.Go
		fmt.Printf("Installing node: %s : %s\n", node.Description, implPath)

		if info.DRY {
			continue
		}
		go func() {
			url := fmt.Sprintf("ipc://%s/%d.sock", info.IPC_ROOT, i)
			var sock mangos.Socket
			var err error
			var msg []byte
			if sock, err = rep.NewSocket(); err != nil {
				die("can't get new rep socket: %s", err)
			}
			sock.AddTransport(ipc.NewTransport())
			sock.AddTransport(tcp.NewTransport())
			if err = sock.Dial(url); err != nil {
				die("can't dial on rep socket: %s", err.Error())
			}
			wg.Add(1)
			for {
				msg, err = sock.Recv()
				cr, err := protocol.fromEncodedMessage(msg)
				if err != nil {
					die("can't decode message: %s", err.Error())
				}
				start := time.Now()
				nresp, err := sandbox.Computes[i](cr)
				end := time.Since(start)
				t = []float{end / time.Second}
				t = append(t, end-t[0])
				var resp []byte
				if err != nil {
					resp, err = protocol.toEncodedMessage(nresp, t)
					if err != nil {
						die("issue encoding your response")
					}
				} else {
					// resp, err = protocol.toErrorBytes(
					if err != nil {
						die()
					}
				}

				err = sock.Send(resp)
				if err != nil {
					die("can't send reply: %s", err.Error())
				}
			}
		}()
		wg.Wait()
	}

	if info.DRY {
		os.Exit(0)
	}
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}
