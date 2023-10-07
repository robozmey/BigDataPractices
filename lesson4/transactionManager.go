package main

import (
	"log"
	"os"
	"sync"
)

var transactionQueue = make(chan string)

var resultQueue = make(chan string)

var transactionMutex sync.Mutex

func handleReplace(newState string) {
	state = newState
	resultQueue <- "OK"
}

func handleGet() {
	resultQueue <- state
}

func handleNGGYU() {
	f, _ := os.Open("nggyu.txt")
	defer f.Close()

	var cache [1000]byte

	n, _ := f.Read(cache[:])

	state = string(cache[:n])
}

type RequestType int32

const (
	GetRequest     RequestType = 0
	ReplaceRequest             = 1
	NGGYURequest               = 2
)

func parseTransaction(request string) (RequestType, string) {
	if request == "get" {
		return GetRequest, ""
	} else if request == "nggyu" {
		return NGGYURequest, ""
	}
	return ReplaceRequest, request[len("replace "):]
}

func transactionHandler() {
	for {
		transaction := <-transactionQueue

		transactionMutex.Lock()

		transactionJournal = append(transactionJournal, transaction)

		log.Println("New transaction came: ", transaction)

		transactionType, rs := parseTransaction(transaction)

		log.Println("New transaction type: ", transactionType)

		switch transactionType {
		case GetRequest:
			handleGet()
		case ReplaceRequest:
			handleReplace(rs)
		case NGGYURequest:
			handleNGGYU()
		}

		log.Println("New transaction handled!")

		transactionMutex.Unlock()
	}
}
