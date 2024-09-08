package akv

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)


func CreateAndrewKeyValueClient(address string) (*AndrewKeyValueClient, error) {
	c, _, err := websocket.DefaultDialer.Dial("ws://"+address+"/ws", nil)
	if err != nil {
		return nil, err
	}
	client := &AndrewKeyValueClient{
		conn: c,
		GetChannel: make(chan Message),
		PutChannel: make(chan Message),
		DeleteChannel: make(chan Message),
		GetLastUpdatedChannel: make(chan Message),
		CacheInvalidationChannel: make(chan Message),
		cache: make(map[Key]Value),
	}
	return client, nil
}

func (client *AndrewKeyValueClient) send(msg Message) (error) {
	err := client.conn.WriteJSON(msg)
    if err != nil {
        return fmt.Errorf("error writing message: %w", err)
    }

	go func() (error) {
		for {
			var response Message
			err = client.conn.ReadJSON(&response)
			if err != nil {
				log.Print("error reading response: %w", err)
				return err
			}

			if (response.Type == INVALIDATE_CACHE) {
				log.Print("Invalidate cache message received");
				client.CacheInvalidationChannel <- response
				continue
			} else if (response.Type == GET_RESPONSE) {
				client.GetChannel <- response
			} else if (response.Type == PUT_RESPONSE) {
				client.PutChannel <- response
			} else if (response.Type == DELETE_RESPONSE) {
				client.DeleteChannel <- response
			} else if (response.Type == GET_LAST_UPDATED_RESPONSE) {
				client.GetLastUpdatedChannel <- response
			}
		}
	}()

	return nil
}

func (client *AndrewKeyValueClient) Get(key string) (Message, error) {
	val, ok := client.TryGetFromCache(key)
	if ok {
		log.Println("Cache hit for key: ", key);
		return Message{Type: GET_RESPONSE, Key: key, Value: val.Value, Timestamp: &val.LastUpdated, Success: true}, nil
	}

	err := client.send(Message{Type: GET, Key: key})
	response := <- client.GetChannel

	if !response.Success {
		return Message{}, errors.New(response.Err);
	}

	client.UpdateCache(key, Value{Value: response.Value, LastUpdated: *response.Timestamp})

	return response, err
}

func (client *AndrewKeyValueClient) Put(key string, value string) (bool, error) {
	err := client.send(Message{Type: PUT, Key: key, Value: value})
	response := <-client.PutChannel

	if !response.Success {
		return false, errors.New(response.Err);
	}

	client.UpdateCache(key, Value{Value: value, LastUpdated: *response.Timestamp})

	return response.Success, err
}

func (client *AndrewKeyValueClient) Delete(key string) (bool, error) {
	err := client.send(Message{Type: DELETE, Key: key})
	response := <-client.DeleteChannel

	if !response.Success {
		return false, errors.New(response.Err);
	}

	client.InvalidateCache(key)

	return response.Success, err
}

func (client *AndrewKeyValueClient) GetLastUpdated(key string) (*time.Time, error) {
	err := client.send(Message{Type: GET_LAST_UPDATED, Key: key})
	response := <-client.GetLastUpdatedChannel

	if !response.Success {
		return nil, errors.New(response.Err);
	}

	return response.Timestamp, err
}

func (client *AndrewKeyValueClient) Close() {
	client.conn.Close()
}

func (client *AndrewKeyValueClient) TryGetFromCache(key string) (Value, bool) {
	value, ok := client.cache[Key(key)]
	return value, ok
}

func (client *AndrewKeyValueClient) InvalidateCache(key string) {
	delete(client.cache, Key(key))
}

func (client *AndrewKeyValueClient) UpdateCache(key string, value Value) {
	client.cache[Key(key)] = value
}

func (client *AndrewKeyValueClient) ListenForCacheInvalidation() {
	for {
		msg := <- client.CacheInvalidationChannel
		log.Printf("Invalidating cache for key: %s", msg.Key)
		client.InvalidateCache(msg.Key)
	}
}