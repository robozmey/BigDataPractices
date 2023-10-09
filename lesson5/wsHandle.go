package main

import (
	"context"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func replicaWriter(c *websocket.Conn, ctx context.Context, peer string) {
	// Sends new transaction
	var transactionSendCount = 0
	var journalSize = len(journal)

	for {
		time.Sleep(100 * time.Millisecond)

		journalSize = len(journal)

		for ; transactionSendCount < journalSize; transactionSendCount++ {
			transaction := journal[transactionSendCount]

			log.Println("WS:", "sending transaction:", transaction, "to:", peer)
			err := wsjson.Write(ctx, c, transaction)
			if err != nil {
				c.Close(websocket.StatusInternalError, err.Error())
				return
			}
		}
	}
}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
	})

	log.Println("WS:", "handled new websocket connection:", r.Host)
	defer log.Println("WS:", "websocket connection handled:", r.Host)

	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	ctx := r.Context()

	var replicationChannel = make(chan Transaction, 100)
	defer close(replicationChannel)

	replicationHandler(c, ctx, r.Host)

	c.Close(websocket.StatusNormalClosure, "")
}
