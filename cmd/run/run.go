package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	mangosock "nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/nano"
	"github.com/redsift/go-sandbox-rpc"
	"github.com/redsift/sandbox-go/sandbox"
)

type result struct {
	response []sandboxrpc.ComputeResponse
	err      map[string]string
}

func main() {
	info, err := sandbox.NewInit(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	panicked := false
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
			var sock nano.Rep
			var err error
			var msg []byte
			if sock, err = mangosock.NewRepSocket(); err != nil {
				die("can't get new rep socket: %s", err)
			}

			if err = sock.Connect(url); err != nil {
				die("can't dial on rep socket: %s", err)
			}

			sendErr := func(nerr error, stack string) {
				resp, err := sandbox.ToErrorBytes(nerr.Error(), stack)
				if err != nil {
					die("issue encoding your error: %s", err)
				}
				_, err = sock.Send(resp)
				if err != nil {
					die("can't send reply: %s", err)
				}
				canSend = false
			}

			defer func() {
				event := recover()
				if event != nil {
					panicked = true
					stack := debug.Stack()
					fmt.Printf("Stack: %s\n", stack)

					err := errors.New("panic")
					if evErr, ok := event.(error); ok {
						err = evErr
					}

					if canSend {
						fmt.Println("canSend")
						sendErr(err, string(stack))
					}
					fmt.Printf("caught a sandbox panic: %s\n", err)
					//die("caught a node panic: %s", err)
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
					sendErr(err, "")
					die("can't decode message: %s", err)
				}
				if _, ok := sandbox.Computes[idx]; !ok {
					sendErr(fmt.Errorf("no node with id: %d", idx), "")
					die("no node with id: %d", idx)
				}
				start := time.Now()
				ch := make(chan *result)
				go func(ch chan<- *result) {
					defer func() {
						event := recover()
						if event != nil {
							stack := debug.Stack()
							fmt.Printf("Stack: %s\n", stack)

							err := errors.New("panic")
							if evErr, ok := event.(error); ok {
								err = evErr
							}

							ch <- &result{
								err: map[string]string{
									"message": err.Error(),
									"stack":   string(stack),
								},
							}

							fmt.Printf("caught a node panic: %s\n", err)
						}
					}()

					nresp, err := sandbox.Computes[idx](cr)
					if err != nil {
						ch <- &result{
							err: map[string]string{
								"message": err.Error(),
							},
						}
						return
					}

					ch <- &result{
						response: nresp,
					}
				}(ch)

				res := <-ch
				if res.err != nil {
					sendErr(errors.New(res.err["message"]), res.err["stack"])
					continue
				}

				end := time.Since(start)
				t := []int64{int64(end / time.Second)}
				t = append(t, int64(end)-t[0])

				resp, err := sandbox.ToEncodedMessage(res.response, t)
				if err != nil {
					sendErr(err, "")
					die("issue encoding your response: %s", err)
				}

				_, err = sock.Send(resp)
				if err != nil {
					die("can't send reply: %s", err)
				}
				canSend = false
			}
		}(url, i)
	}
	wg.Wait()

	if panicked {
		select {} // wait to get killed
	}
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}
