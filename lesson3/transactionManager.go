package main

import (
	"fmt"
	"os"
)

var transactionQueue = make(chan string)

var resultQueue = make(chan string)

func handleReplace(newData string) {
	data = newData
	resultQueue <- "OK"
}

func handleGet() {
	resultQueue <- data
}

func handleNGGYU() {
	f, _ := os.Open("nggyu.txt")
	defer f.Close()

	var cache [1000]byte

	n, _ := f.Read(cache[:])

	data = string(cache[:n])
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

		journal = append(journal, transaction)

		fmt.Println("New transaction came: ", transaction)

		transactionType, rs := parseTransaction(transaction)

		fmt.Println("New transaction type: ", transactionType)

		switch transactionType {
		case GetRequest:
			handleGet()
		case ReplaceRequest:
			handleReplace(rs)
		case NGGYURequest:
			handleNGGYU()
		}

		fmt.Println("New transaction handled!")
	}
}
