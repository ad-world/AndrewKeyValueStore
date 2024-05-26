package akv

import (
	"errors"
	"log"
	"time"
)

type AndrewKeyValueStore struct {
	Store map[Key]Value
}

// CreateAndrewKeyValueStore creates a new instance of AndrewKeyValueStore with an empty store and returns a pointer to it.
func CreateAndrewKeyValueStore() *AndrewKeyValueStore {
	return &AndrewKeyValueStore{
		Store: make(map[Key]Value),
	}
}

// Get retrieves the value of a key from the store. 
// It will populate the reply with the value of the key if it exists, otherwise it will return an error.
func (store *AndrewKeyValueStore) Get(args *GetRequest, reply *string) error {
	// TODO: add locking mechanism before reading the store
	log.Println("Get request for key ", args.Key)
	value, ok := store.Store[Key(args.Key)]
	if !ok {
		*reply = ""
		log.Println("Key " + args.Key + " not found")
		return errors.New("Key " + args.Key + " not found")
	}

	log.Println("Value for key ", args.Key, " found.")
	// TODO: unlock the store

	*reply = value.Value

	return nil
}

// Put inserts a key-value pair into the store. 
// It will populate the reply with true if the operation is successful, otherwise it will return an error.
func (store *AndrewKeyValueStore) Put(args *PutRequest, reply *bool) error {
	log.Println("Put request for key ", args.Key)
	// TODO: add locking mechanism before updating the store
	store.Store[Key(args.Key)] = Value{
		Value: args.Value,
		LastUpdated: time.Now(),
	}

	// TODO: unlock the store

	*reply = true
	return nil
}

// Delete removes a key from the store. 
// It will populate the reply with true if the operation is successful, otherwise it will return an error.
func (store *AndrewKeyValueStore) Delete(args *DeleteRequest, reply *bool) error {
	log.Println("Delete request for key ", args.Key)
	// TODO: add locking mechanism before updating the store
	_, ok := store.Store[Key(args.Key)]
	
	if ok {
		delete((store.Store), Key(args.Key))
		log.Println("Key " + args.Key + " deleted.")
		*reply = true
	} else {
		log.Println("Key " + args.Key + " not found.")
		*reply = false
	}
	// TODO: unlock the store

	return nil
}

// GetLastUpdated retrieves the last updated time of a key from the store. 
// It will populate the reply with the last updated time of the key if it exists, otherwise it will return an error.
func (store *AndrewKeyValueStore) GetLastUpdated(args *GetLastUpdatedRequest, reply *time.Time) error {
	log.Println("GetLastUpdated request for key ", args.Key)
	item, ok := store.Store[Key(args.Key)]

	if !ok {
		log.Println("Key " + args.Key + " not found")
		return errors.New("Key " + args.Key + " not found")
	}

	log.Println("Last updated time for key ", args.Key, " found.")
	*reply = item.LastUpdated

	return nil
}