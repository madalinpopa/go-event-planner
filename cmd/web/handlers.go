package main

import "net/http"

func ping(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		panic(err)
	}
}

func home(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("home"))
	if err != nil {
		panic(err)
	}
}
