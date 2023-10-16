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

	log.Println("TM:", "Got new transaction: ", transaction)

	transactionMutex.Lock()
	defer transactionMutex.Unlock()

	var isTransactionApplied = vclock[transaction.Source] >= transaction.Id
	if isTransactionApplied {
		resultQueue <- http.StatusOK
		log.Println("TM:", "Transaction already applied: ", transaction.Source, vclock[transaction.Source])
		return
	}

	vclock[transaction.Source] = transaction.Id

	patch, err := jsonpatch.DecodePatch([]byte(transaction.Payload))
	if err != nil {
		resultQueue <- http.StatusBadRequest
		log.Println("TM:", "Cannot decode transaction: ", err)
		return
	}

	journal = append(journal, transaction)

	newsnap, err := patch.Apply([]byte(snap))
	if err != nil {
		resultQueue <- http.StatusBadRequest
		log.Println("TM:", "Cannot apply transaction: ", err)
		return
	}
	snap = string(newsnap)

	resultQueue <- http.StatusOK
	log.Println("TM:", "Transaction successfully applied")
}

func transactionHandler() {
	for {
		handleTransaction()
	}
}
