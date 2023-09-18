package main

import (
	"log"
	"net/http"
)

func requestReplace(w http.ResponseWriter, r *http.Request) {
	println("replace")
}

func requestGet(w http.ResponseWriter, r *http.Request) {
	println("get")
}

func main() {
	http.HandleFunc("/replace", requestReplace) // Устанавливаем роутер
	http.HandleFunc("/get", requestGet)
	err := http.ListenAndServe(":8080", nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
