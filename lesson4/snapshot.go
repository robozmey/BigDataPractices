package main

import (
	"log"
	"time"
)

type Snapshot struct {
	journal Journal
	state   string
}

var snapshot = Snapshot{make(Journal, 0), ""}

func snapshotHandler() {
	for {
		transactionMutex.Lock()
		log.Println("Saving snapshot...")
		snapshot = Snapshot{transactionJournal, state}
		log.Println("Snapshot saved:", snapshot)
		transactionMutex.Unlock()

		time.Sleep(5 * time.Second)
	}
}
