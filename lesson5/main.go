package main

import (
	"fmt"
	"os"
)

func main() {
	argCount := len(os.Args)

	var port = "8080"

	if argCount >= 2 {
		port = os.Args[1]
	}

	if argCount >= 3 {
		mySource = os.Args[2]
	}

	go httpServer(port)
	go transactionHandler()
	go replication()

	fmt.Scanln()
}
