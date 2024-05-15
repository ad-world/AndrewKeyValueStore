package akv

import (
	"errors"
	"time"
)

type AndrewKeyValueStore struct {
	Store map[Key]Value
}

func CreateAndrewKeyValueStore() *AndrewKeyValueStore {
	return &AndrewKeyValueStore{
		Store: make(map[Key]Value),
	}
}

func (store *AndrewKeyValueStore) Get(args *GetRequest, reply *string) error {
	// TODO: add locking mechanism before reading the store
	value, ok := store.Store[Key(args.Key)]
	if !ok {
		*reply = ""
		return errors.New("Key " + args.Key + "  not found")
	}

	// TODO: unlock the store

	*reply = value.Value

	return nil
}

func (store *AndrewKeyValueStore) Put(args *PutRequest, reply *bool) error {
	// TODO: add locking mechanism before updating the store
	store.Store[Key(args.Key)] = Value{
		Value: args.Value,
		LastUpdated: time.Now(),
	}

	// TODO: unlock the store

	*reply = true
	return nil
}

func (store *AndrewKeyValueStore) Delete(args *DeleteRequest, reply *bool) error {
	// TODO: add locking mechanism before updating the store
	_, ok := store.Store[Key(args.Key)]
	
	if ok {
		delete((store.Store), Key(args.Key))
		*reply = true
	} else {
		*reply = false
	}
	// TODO: unlock the store

	return nil
}