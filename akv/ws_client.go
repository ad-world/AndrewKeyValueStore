package akv

import (
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
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
}

func CreateAndrewKeyValueClient(address string) (*AndrewKeyValueClient, error) {
	c, _, err := websocket.DefaultDialer.Dial("ws://"+address+"/ws", nil)
	if err != nil {
		return nil, err
	}
	return &AndrewKeyValueClient{conn: c}, nil
}

func (client *AndrewKeyValueClient) sendAndReceive(msg Message) (Message, error) {
	err := client.conn.WriteJSON(msg)
    if err != nil {
        return Message{}, fmt.Errorf("error writing message: %w", err)
    }

    var response Message
    err = client.conn.ReadJSON(&response)
    if err != nil {
        return Message{}, fmt.Errorf("error reading response: %w", err)
    }

	if (response.Type == ERROR) {
		return Message{}, errors.New(response.Err)
	}

    return response, nil
}

func (client *AndrewKeyValueClient) Get(key string) (Message, error) {
	response, err := client.sendAndReceive(Message{Type: GET, Key: key})
	return response, err
}

func (client *AndrewKeyValueClient) Put(key string, value string) (bool, error) {
	response, err := client.sendAndReceive(Message{Type: PUT, Key: key, Value: value})
	return response.Success, err
}

func (client *AndrewKeyValueClient) Delete(key string) (bool, error) {
	response, err := client.sendAndReceive(Message{Type: DELETE, Key: key})
	return response.Success, err
}

func (client *AndrewKeyValueClient) GetLastUpdated(key string) (*time.Time, error) {
	response, err := client.sendAndReceive(Message{Type: GET_LAST_UPDATED, Key: key})
	return response.Timestamp, err
}

func (client *AndrewKeyValueClient) Close() {
	client.conn.Close()
}
