package akv

import (
	"time"

	"github.com/gorilla/websocket"
)

type Key string
type Value struct {
	Value string 		     `json:"value"`
	LastUpdated time.Time	 `json:"last_updated"`
}

type DeleteRequest struct {
	Key string
}

type GetRequest struct {
	Key string
}

type PutRequest struct {
	Key string
	Value string
}

type GetLastUpdatedRequest struct {
	Key string
}


type AndrewKeyValueStore struct {
	Store map[Key]Value
}

type MessageType int

const (
	GET MessageType = 0
	PUT MessageType = 1
	DELETE MessageType = 2
	GET_LAST_UPDATED MessageType = 3
	GET_RESPONSE MessageType = 4
	PUT_RESPONSE MessageType = 5
	DELETE_RESPONSE MessageType = 6
	GET_LAST_UPDATED_RESPONSE MessageType = 7
	ERROR MessageType = 8
	INVALIDATE_CACHE MessageType = 9
)

type Message struct {
	Type      MessageType  `json:"type"`
	Key       string  `json:"key"`
	Value     string  `json:"value"`
	Timestamp *time.Time  `json:"timestamp"`
	Success   bool   `json:"success"`
	Err       string  `json:"err"`
}

type AndrewKeyValueClient struct {
	conn *websocket.Conn
	GetChannel chan Message
	PutChannel chan Message
	DeleteChannel chan Message
	GetLastUpdatedChannel chan Message
	CacheInvalidationChannel chan Message
	cache map[Key]Value
}
