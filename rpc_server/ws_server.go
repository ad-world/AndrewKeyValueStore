package main

import (
	"akv/akv"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(ws *websocket.Conn, store *akv.AndrewKeyValueStore) {
	defer ws.Close();

	for {
		var msg akv.Message
		err := ws.ReadJSON(&msg)

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("Client closed connection normally")
			} else if err == io.EOF || websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				log.Println("Client disconnected unexpectedly")
			} else {
				log.Printf("Error reading message: %v", err)
			}
			return // Exit this goroutine, but keep the server running
		}

		switch msg.Type {
			case akv.GET:
				log.Println("GET request for Key: ", msg.Key);

				value, err := store.Get(&akv.GetRequest{Key: msg.Key})

				if err != nil {
					log.Printf("Error in GET: %v", err)
					response := akv.Message{Type: akv.ERROR, Err: err.Error()}
					ws.WriteJSON(response)
					continue
				}

				response := akv.Message{Type: akv.GET_RESPONSE, Key: msg.Key, Value: value.Value, Timestamp: &value.LastUpdated}
				ws.WriteJSON(response)
			case akv.PUT:
				log.Println("PUT request for Key: ", msg.Key);

				success, err := store.Put(&akv.PutRequest{Key: msg.Key, Value: msg.Value})
				if err != nil {
					log.Printf("Error in PUT: %v", err)
					response := akv.Message{Type: akv.ERROR, Err: err.Error()}
					ws.WriteJSON(response)
					continue
				}

				response := akv.Message{Type: akv.PUT_RESPONSE, Key: msg.Key, Value: msg.Value, Success: success}
				ws.WriteJSON(response)
			case akv.DELETE:
				log.Println("DELETE request for Key: ", msg.Key);

				success, err := store.Delete(&akv.DeleteRequest{Key: msg.Key})
				if err != nil {
					log.Printf("Error in DELETE: %v", err)
					response := akv.Message{Type: akv.ERROR, Err: err.Error()}
					ws.WriteJSON(response)
					continue
				}
				response := akv.Message{Type: akv.DELETE_RESPONSE, Key: msg.Key, Success: success}
				ws.WriteJSON(response)
			case akv.GET_LAST_UPDATED:
				log.Println("GET_LAST_UPDATED request for Key: ", msg.Key);

				timestamp, err := store.GetLastUpdated(&akv.GetLastUpdatedRequest{Key: msg.Key})
				if err != nil {
					log.Printf("Error in GET_LAST_UPDATEd: %v", err)
					response := akv.Message{Type: akv.ERROR, Err: err.Error() }
					ws.WriteJSON(response)
					continue
				}
				response := akv.Message{Type: akv.GET_LAST_UPDATED_RESPONSE, Key: msg.Key, Timestamp: timestamp } 
				ws.WriteJSON(response)
		}
	}
}

func handleConnections(store *akv.AndrewKeyValueStore, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }
	go handleWebSocket(ws, store);
}

func main() {
	store := akv.CreateAndrewKeyValueStore()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(store, w, r)
	})

	fmt.Println("Server started on port 1234")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
