package main

import (
	"sandbox-go/sandbox"
	"os"
	"fmt"
)

func main(){
	_, err := sandbox.NewInit(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
