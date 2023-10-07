package main

import (
	"context"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"os"
	"strings"
	"time"
)

func replicationDial(peer string) {
	//ctx, cancel := context.WithTimeout(, time.Second*10)
	//defer cancel()
	var ctx = context.Background()

	var url = fmt.Sprintf("ws://%s/ws", peer)
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		log.Println("Repl:", "Unable to connect", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	for {
		log.Println("Repl:", "Waiting transactions from", peer)
		var transaction Transaction
		//_, b, err := c.Read(ctx)
		//log.Println(string(b))
		err := wsjson.Read(ctx, c, &transaction)
		if err != nil {
			log.Println("Repl:", "Stopped waiting transactions from", peer, "because:", err)
			break
		}
		log.Println("Repl:", "From", peer, "got transaction", transaction)
		transactionQueue <- transaction
		<-resultQueue
	}

	c.Close(websocket.StatusNormalClosure, "")
}

func replicationRoutine(peer string) {
	for {
		replicationDial(peer)
		time.Sleep(1 * time.Second)
	}
}

func replication() {
	f, _ := os.Open("peers.txt")
	defer f.Close()

	var cache [1000]byte

	n, _ := f.Read(cache[:])

	peers = strings.Split(string(cache[:n]), "\n")

	for _, peer := range peers {
		go replicationRoutine(peer)
	}

}
