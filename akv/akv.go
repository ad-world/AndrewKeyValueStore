package akv

import (
	"errors"
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
func (store *AndrewKeyValueStore) Get(args *GetRequest) (*Value, error) {
	// TODO: add locking mechanism before reading the store
	value, ok := store.Store[Key(args.Key)]
	if !ok {
		return nil, errors.New("Key '" + args.Key + "' not found")
	}

	// TODO: unlock the store

	return &value, nil
}

// Put inserts a key-value pair into the store. 
// It will populate the reply with true if the operation is successful, otherwise it will return an error.
func (store *AndrewKeyValueStore) Put(args *PutRequest) (bool, error) {
	// TODO: add locking mechanism before updating the store
	store.Store[Key(args.Key)] = Value{
		Value: args.Value,
		LastUpdated: time.Now(),
	}

	// TODO: unlock the store
	return true, nil
}

// Delete removes a key from the store. 
// It will populate the reply with true if the operation is successful, otherwise it will return an error.
func (store *AndrewKeyValueStore) Delete(args *DeleteRequest) (bool, error) {
	// TODO: add locking mechanism before updating the store
	_, ok := store.Store[Key(args.Key)]
	
	if ok {
		delete((store.Store), Key(args.Key))
		return true, nil
	} else {
		return false, errors.New("Key '" + args.Key + "' not found")
	}
	// TODO: unlock the store
}

// GetLastUpdated retrieves the last updated time of a key from the store. 
// It will populate the reply with the last updated time of the key if it exists, otherwise it will return an error.
func (store *AndrewKeyValueStore) GetLastUpdated(args *GetLastUpdatedRequest) (*time.Time, error) {
	item, ok := store.Store[Key(args.Key)]

	if !ok {
		return nil, errors.New("Key '" + args.Key + "' not found")
	}

	return &item.LastUpdated, nil
}