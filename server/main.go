package main

import (
	"akv/akv"
	"akv/ws_server"
	"fmt"
	"log"
	"net/http"
)


func main() {
	store := akv.CreateAndrewKeyValueStore()
	server := ws_server.CreateWsServer()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.HandleConnections(store, w, r)
	})

	go server.HandleBroadcast()

	fmt.Println("Server started on port 1234")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
