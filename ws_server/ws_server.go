package ws_server

import (
	"akv/akv"
	"io"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)


func CreateWsServer() *WsServer {
	return &WsServer{
		clients: make(map[*websocket.Conn]*ClientInfo),
		broadcast: make(chan akv.Message),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (server *WsServer) BroadcastMessage(msg akv.Message) {
	server.broadcast <- msg
}

func (server *WsServer) HandleBroadcast() {
	for msg := range server.broadcast {
		for _, client := range server.clients {
			select {
			case client.send <- msg:
			default: 
				close(client.send)
				delete(server.clients, client.conn)
			}
		}
	}
}

func (server *WsServer) HandleClient(client *ClientInfo) {
	defer func() {
		client.conn.Close();
		delete(server.clients, client.conn);
	}()

	for msg := range client.send {
		err := client.conn.WriteJSON(msg);
		if err != nil {
			log.Printf("Error sending message to client with IP %s", client.conn.RemoteAddr())
			return
		}
	}
}

func (server* WsServer) HandleWebsocket(ws *websocket.Conn, store *akv.AndrewKeyValueStore) {
	client := &ClientInfo{
		conn: ws,
		send: make(chan akv.Message, 256),
	}

	server.clients[ws] = client

	go server.HandleClient(client)

	defer func() {
		delete(server.clients, ws)
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
				
				server.BroadcastMessage(akv.Message{Type: akv.INVALIDATE_CACHE, Key: msg.Key});
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

func (server *WsServer) HandleConnections(store *akv.AndrewKeyValueStore, w http.ResponseWriter, r *http.Request) {
	ws, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }
	go server.HandleWebsocket(ws, store);
}