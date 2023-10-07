package main

import (
	"context"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func wsHandle(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
	})

	log.Println("WS:", "handled new websocket connection")
	defer log.Println("WS:", "websocket connection handled")

	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	var replicationChannel = make(chan Transaction, 100)
	defer close(replicationChannel)

	{
		transactionMutex.Lock()

		// Register channel
		replicationChannels = append(replicationChannels, replicationChannel)

		// Sends all journaled transactions
		for _, transaction := range journal {
			log.Println("WS:", "sending transaction:", transaction)
			wsjson.Write(ctx, c, transaction)
		}

		transactionMutex.Unlock()
	}

	// Sends new transaction from channel
	for transaction := range replicationChannel {
		log.Println("WS:", "sending transaction:", transaction)
		wsjson.Write(ctx, c, transaction)
	}

	c.Close(websocket.StatusNormalClosure, "")
}
