package main

import (
	"embed"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
)

//go:embed index.html
var content embed.FS

func getHandle(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(snap))

	log.Println("handled get, send snap: %s", snap)
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
}

func vclockHandle(w http.ResponseWriter, r *http.Request) {

	bytes, err := json.Marshal(vclock)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)

	log.Printf("handled vclock, send vclock: %s", bytes)
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
