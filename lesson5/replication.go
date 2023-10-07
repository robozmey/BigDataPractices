package main

import (
	"context"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"os"
	"strings"
)

func replicationRoutine(peer string) {
	//ctx, cancel := context.WithTimeout(, time.Second*10)
	//defer cancel()
	var ctx = context.Background()

	var url = fmt.Sprintf("ws://%s/ws", peer)
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		log.Println("Unable to connect", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	for {
		log.Println("Waiting for transaction from ", peer)
		var transaction Transaction
		//_, b, err := c.Read(ctx)
		//log.Println(string(b))
		err := wsjson.Read(ctx, c, &transaction)
		if err != nil {
			log.Println("Stopped waiting for transaction from ", peer, "because", err)
			break
		}
		log.Println("From ", peer, "got transaction", transaction)
		transactionQueue <- transaction
	}

	c.Close(websocket.StatusNormalClosure, "")
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
