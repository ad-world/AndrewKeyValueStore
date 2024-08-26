package ws_server

import (
	"akv/akv"
	"github.com/gorilla/websocket"
)

type ClientInfo struct {
	conn *websocket.Conn
	send chan akv.Message
}

type WsServer struct {
	clients map[*websocket.Conn]*ClientInfo
	broadcast chan akv.Message
	upgrader websocket.Upgrader
}
