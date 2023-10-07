package main

import (
	"log/slog"
	"net/http"
)

func requestReplace(w http.ResponseWriter, r *http.Request) {
	var buff [1000]byte

	n, _ := r.Body.Read(buff[:])
	defer r.Body.Close()

	transactionQueue <- "replace " + string(buff[:n])

	result := <-resultQueue

	if result == "OK" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func requestGet(w http.ResponseWriter, r *http.Request) {
	transactionQueue <- "get"

	result := <-resultQueue

	w.Write([]byte(result))
}

func requestNGGYU(w http.ResponseWriter, r *http.Request) {
	transactionQueue <- "nggyu"
}

func httpServer() {
	http.HandleFunc("/replace", requestReplace) // Устанавливаем роутер
	http.HandleFunc("/get", requestGet)
	http.HandleFunc("/NGGYU", requestNGGYU)

	slog.Info("Server started...")
	http.ListenAndServe(":8080", nil) // ус
}
