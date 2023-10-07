package main

import (
	"context"
	"embed"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

//go:embed index.html
var content embed.FS

func getHandle(w http.ResponseWriter, r *http.Request) {
	//c, err := websocket.Accept(w, r, nil)
	//if err != nil {
	//	panic(err)
	//}
	//defer c.Close(websocket.StatusInternalError, "the sky is falling")
	//
	//ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	//defer cancel()

	w.Write([]byte(snap))

	//err = c.Write(ctx, websocket.MessageText, []byte(snap))
	//if err != nil {
	//	panic(err)
	//}
	//
	log.Printf("handled get, send snap: %s", snap)
	//
	//c.Close(websocket.StatusNormalClosure, "")
}

var transactionCount uint64 = 1

func postHandle(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var transactionId = transactionCount

	transactionCount++

	var transaction = Transaction{Source: mySource, Id: transactionId, Payload: string(body)}

	transactionQueue <- transaction
	result := <-resultQueue

	w.WriteHeader(result)

	//c, err := websocket.Accept(w, r, nil)
	//if err != nil {
	//	panic(err)
	//}
	//defer c.Close(websocket.StatusInternalError, "the sky is falling")
	//
	//ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	//defer cancel()
	//
	//var v interface{}
	//err = wsjson.Read(ctx, c, &v)
	//if err != nil {
	//	// ...
	//}
	//
	//log.Printf("received: %v", v)
	//
	//c.Close(websocket.StatusNormalClosure, "")
}

func vclockHandle(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(vclock)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)

	log.Printf("handled vclock, send vclock: %s", bytes)
}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
	})

	log.Printf("handled new websocket connection")

	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	transactionMutex.Lock()
	defer transactionMutex.Unlock()

	for _, transaction := range journal {

		log.Println("sending transaction:", transaction)
		wsjson.Write(ctx, c, transaction)

	}

	log.Println("websocket connection handled")

	c.Close(websocket.StatusNormalClosure, "")
}

func httpServer(port string) {
	slog.Info(port)
	http.Handle("/test/", http.StripPrefix("/test/", http.FileServer(http.FS(content))))
	http.HandleFunc("/vclock", vclockHandle)
	http.HandleFunc("/post", postHandle)
	http.HandleFunc("/get", getHandle)
	http.HandleFunc("/ws", wsHandle)

	slog.Info("Server started...")
	http.ListenAndServe(":"+port, nil) // ус
}
