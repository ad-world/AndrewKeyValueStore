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

type ClientInfo struct {
	conn *websocket.Conn
	send chan akv.Message
}

// Create a map of connected clients
var clients = make(map[*websocket.Conn]*ClientInfo)

// Create a channel of message type
var broadcast = make(chan akv.Message)

// Function that sends a message to the channel
func broadcastMessage(msg akv.Message) {
    broadcast <- msg
}

func handleBroadcasts() {
	// Whenever a message is send to the channel
	for msg := range broadcast {
        for _, client := range clients {
            select {
                // Message sent to client's channel
            case client.send <- msg:
            default:
                // Client's channel is full, close the connection
                close(client.send)
                delete(clients, client.conn)
            }
        }
    }
}

func handleClient(client *ClientInfo) {
	defer func() {
		client.conn.Close();
		delete(clients, client.conn)
	}()

	for msg := range client.send {
		log.Printf("Sending message to client");
		err := client.conn.WriteJSON(msg);
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
            return
		}
	}
}

func handleWebSocket(ws *websocket.Conn, store *akv.AndrewKeyValueStore) {
	client := &ClientInfo{
		conn: ws,
		send: make(chan akv.Message, 256),
	}

	clients[ws] = client

	go handleClient(client);

	defer func() {
		delete(clients, ws)
		close(client.send)
	}()

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
					response := akv.Message{Type: akv.GET_RESPONSE, Err: err.Error(), Success: false}
					client.send <- response
					continue
				}

				response := akv.Message{Type: akv.GET_RESPONSE, Key: msg.Key, Value: value.Value, Timestamp: &value.LastUpdated, Success: true}
				client.send <- response
			case akv.PUT:
				log.Println("PUT request for Key: ", msg.Key);

				success, err := store.Put(&akv.PutRequest{Key: msg.Key, Value: msg.Value})
				if err != nil {
					log.Printf("Error in PUT: %v", err)
					response := akv.Message{Type: akv.PUT_RESPONSE, Err: err.Error(), Success: false}
					client.send <- response
					continue
				}
				
				response := akv.Message{Type: akv.PUT_RESPONSE, Key: msg.Key, Value: msg.Value, Success: success}
				client.send <- response
				
				broadcastMessage(akv.Message{Type: akv.INVALIDATE_CACHE, Key: msg.Key});
			case akv.DELETE:
				log.Println("DELETE request for Key: ", msg.Key);

				success, err := store.Delete(&akv.DeleteRequest{Key: msg.Key})
				if err != nil {
					log.Printf("Error in DELETE: %v", err)
					response := akv.Message{Type: akv.DELETE_RESPONSE, Err: err.Error(), Success: false}
					client.send <- response
					continue
				}
				response := akv.Message{Type: akv.DELETE_RESPONSE, Key: msg.Key, Success: success}
				client.send <- response
			case akv.GET_LAST_UPDATED:
				log.Println("GET_LAST_UPDATED request for Key: ", msg.Key);

				timestamp, err := store.GetLastUpdated(&akv.GetLastUpdatedRequest{Key: msg.Key})
				if err != nil {
					log.Printf("Error in GET_LAST_UPDATEd: %v", err)
					response := akv.Message{Type: akv.GET_LAST_UPDATED_RESPONSE, Err: err.Error(), Success: false }
					client.send <- response
					continue
				}
				response := akv.Message{Type: akv.GET_LAST_UPDATED_RESPONSE, Key: msg.Key, Timestamp: timestamp, Success: true } 
				client.send <- response
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

	go handleBroadcasts()

	fmt.Println("Server started on port 1234")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
