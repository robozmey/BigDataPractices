package main

import (
	"log/slog"
	"net/http"
	"os"
)

func writeCacheFile(cache []byte) {
	f, _ := os.Create("cache.txt")
	defer f.Close()

	f.Write(cache)
}

func readCacheFile(cache []byte) {
	f, _ := os.Open("cache.txt")
	defer f.Close()

	f.Read(cache)
}

func requestReplace(w http.ResponseWriter, r *http.Request) {
	println("replace")

	var cache [1000]byte

	r.Body.Read(cache[:])
	defer r.Body.Close()

	println(string(cache[:]))

	writeCacheFile(cache[:])

	w.WriteHeader(http.StatusOK)
}

func requestGet(w http.ResponseWriter, r *http.Request) {
	println("get")

	var cache [1000]byte

	readCacheFile(cache[:])

	println(string(cache[:]))

	w.Write(cache[:])
}

func requestNGGYU(w http.ResponseWriter, r *http.Request) {
	println("nggyu")

	nggyu, _ := os.ReadFile("nggyu.txt")

	writeCacheFile(nggyu)
}

func main() {
	http.HandleFunc("/replace", requestReplace) // Устанавливаем роутер
	http.HandleFunc("/get", requestGet)
	http.HandleFunc("/NGGYU", requestNGGYU)

	slog.Info("Server started...")
	http.ListenAndServe(":8080", nil) // устанавливаем порт веб-сервера

}
