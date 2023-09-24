package main

func main() {
	go httpServer()
	go transactionHandler()

	for {
	}
}
