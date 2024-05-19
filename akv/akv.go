package akv

import (
	"errors"
	"log"
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