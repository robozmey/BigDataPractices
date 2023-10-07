package main

import (
	jsonpatch "github.com/evanphx/json-patch/v5"
	"log"
	"net/http"
	"sync"
)

var transactionQueue = make(chan Transaction)
var resultQueue = make(chan int)
var transactionMutex sync.Mutex

func handleTransaction() {
	transaction := <-transactionQueue

	log.Println("Got new transaction: ", transaction)

	transactionMutex.Lock()
	defer transactionMutex.Unlock()

	var isTransactionApplied = vclock[transaction.Source] >= transaction.Id

	if isTransactionApplied {
		resultQueue <- http.StatusOK
		log.Println("Transaction already applied: ", transaction.Source, vclock[transaction.Source])
		return
	}

	vclock[transaction.Source] = transaction.Id

	journal = append(journal, transaction)

	patch, err := jsonpatch.DecodePatch([]byte(transaction.Payload))
	if err != nil {
		resultQueue <- http.StatusBadRequest
		log.Println("Cannot decode transaction: ", err)
		return
	}

	newsnap, err := patch.Apply([]byte(snap))
	if err != nil {
		resultQueue <- http.StatusBadRequest
		log.Println("Cannot apply transaction: ", err)
		return
	}
	snap = string(newsnap)

	resultQueue <- http.StatusOK
	log.Println("Transaction successfully applied")
}

func transactionHandler() {
	for {
		handleTransaction()
	}
}
