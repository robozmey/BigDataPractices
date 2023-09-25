package main

import (
	"fmt"
)

func main() {
	go httpServer()
	go transactionHandler()
	go snapshotHandler()

	fmt.Scanln()
}
