package main

import (
	"fmt"
	"os"
)

func main() {
	argCount := len(os.Args)

	if argCount <= 1 {
		go httpServer("8080")
	} else {
		go httpServer(os.Args[1])
	}
	go transactionHandler()
	go replication()

	fmt.Scanln()
}
